// dnsilly - dns automation utility
// Copyright (C) 2025  bitrate16 (bitrate16@gmail.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"dnsilly/config"
	"dnsilly/rules"
	"dnsilly/server"
	"dnsilly/triggers"
	"dnsilly/util"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	configPath := config.GetConfigPath()

	// Used to signal worker thread to stop
	onExit := make(chan os.Signal, 1)
	// Used to wait for worker thread stop
	onExited := make(chan struct{})
	// Handle error from server
	onError := make(chan struct{})

	// Status code
	exitCode := 0

	// Main logic worker thread
	go func() {
		isRunning := true

	workerLoop:
		for isRunning {
			// Prepare config
			fmt.Printf("[%s] Loading config: %s\n", util.Now(), configPath)
			conf, err := config.ParseConfig(configPath)
			if err != nil {
				fmt.Printf("[%s] Error while reading config: %v\n", util.Now(), err)
				isRunning = false

				// Hard fail
				exitCode = 1
				break workerLoop
			}
			configModTime, err := util.GetFileModificationTime(configPath)
			if err != nil {
				fmt.Printf("[%s] Error while checking config modification time: %v\n", util.Now(), err)
				isRunning = false

				// Hard fail
				exitCode = 1
				break workerLoop
			}

			// Prepare rules
			fmt.Printf("[%s] Loading rules: %s\n", util.Now(), conf.Rules)
			dnsRules, err := rules.ParseRules(conf.Rules)
			if err != nil {
				fmt.Printf("[%s] Error while reading rules: %s\n", util.Now(), err)

				// This is not fail, just do nothing
				dnsRules = nil
			}
			rulesModTime, err := util.GetFileModificationTime(conf.Rules)
			if err != nil {
				fmt.Printf("[%s] Error while checking rules modification time: %v\n", util.Now(), err)

				// This is not fail, just do nothing
				dnsRules = nil
			}

			// Start server
			dnsServer := server.NewServer(conf)
			dnsServer.SetRules(dnsRules)
			go func() {
				err = dnsServer.Start()
				if err != nil {
					fmt.Printf("[%s] Error while starting server: %v\n", util.Now(), err)
					onError <- struct{}{}
				}
			}()

			// Periodically check for config updates
		checkLoop:
			for {
				// No autoreload
				if conf.Reload == 0 {
					select {
					case signal := <-onExit:
						fmt.Printf("Handler signal %v\n", signal)
						isRunning = false

						break checkLoop
					case <-onError:
						isRunning = false

						// Trigger lifecycle
						triggers.TriggerLifecycle(conf, triggers.OnStop)

						// Hard fail
						exitCode = 1
						break workerLoop
					}
				} else { // Autoreload
					timer := time.NewTimer(conf.Reload)

					select {
					case signal := <-onExit:
						fmt.Printf("Handler signal %v\n", signal)
						isRunning = false

						break checkLoop
					case <-timer.C:
						// Check config modification time
						newConfigModTime, err := util.GetFileModificationTime(configPath)
						if err != nil {
							fmt.Printf("[%s] Error while checking config modification time: %v\n", util.Now(), err)

							// Ignore and skip iteration
							continue checkLoop
						}

						// Update config
						if configModTime != newConfigModTime {
							fmt.Printf("[%s] Reloading config: %s\n", util.Now(), configPath)

							newConf, err := config.ParseConfig(configPath)
							if err != nil {
								fmt.Printf("[%s] Error while reading config: %v\n", util.Now(), err)

								// Hard fail
								exitCode = 1
								break workerLoop
							}

							// Stop server
							err = dnsServer.Stop()
							if err != nil {
								fmt.Printf("[%s] Error while stopping server: %v\n", util.Now(), err)
								isRunning = false

								// Hard fail
								exitCode = 1
								break workerLoop
							}

							// Trigger lifecycle
							triggers.TriggerLifecycle(conf, triggers.OnPartialStop)

							conf = newConf

							// Check rules befire starting server
							// ...
						}

						// Check rules modification time
						newRulesModTime, err := util.GetFileModificationTime(conf.Rules)
						if err != nil {
							fmt.Printf("[%s] Error while checking rules modification time: %v\n", util.Now(), err)

							// This is not fail, just do nothing
							dnsRules = nil
						}

						// Update rules
						if (newConfigModTime != configModTime) || (newRulesModTime != rulesModTime) {
							fmt.Printf("[%s] Reloading rules: %s\n", util.Now(), conf.Rules)

							newRules, err := rules.ParseRules(conf.Rules)
							if err != nil {
								fmt.Printf("[%s] Error while reading rules: %s\n", util.Now(), err)

								// This is not fail, just do nothing
								dnsRules = nil
							}

							dnsRules = newRules
							rulesModTime = newRulesModTime

							dnsServer.SetRules(dnsRules)
						}

						// Continue with server restart
						if newConfigModTime != configModTime {
							configModTime = newConfigModTime

							// Start server
							dnsServer = server.NewServer(conf)
							dnsServer.SetRules(dnsRules)
							go func() {
								err = dnsServer.Start()
								if err != nil {
									fmt.Printf("[%s] Error while starting server: %v\n", util.Now(), err)
									onError <- struct{}{}
								}
							}()

							// Trigger lifecycle
							triggers.TriggerLifecycle(conf, triggers.OnPartialStart)
						}

					case <-onError:
						isRunning = false

						// Trigger lifecycle
						triggers.TriggerLifecycle(conf, triggers.OnStop)

						// Hard fail
						exitCode = 1
						break workerLoop
					}
				}
			}

			// Stop server
			err = dnsServer.Stop()
			if err != nil {
				fmt.Printf("[%s] Error while stopping server: %v\n", util.Now(), err)
				isRunning = false

				// Hard fail
				exitCode = 1
				break workerLoop
			}

			// Trigger lifecycle
			triggers.TriggerLifecycle(conf, triggers.OnStop)
		}

		onExited <- struct{}{}
	}()

	// Handle signals
	signal.Notify(onExit, os.Interrupt)

	fmt.Printf("[%s] %s", util.Now(), "Service starting\n")

	// Wait for exit
	<-onExited

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
