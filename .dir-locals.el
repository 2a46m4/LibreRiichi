;;; Directory Local Variables            -*- no-byte-compile: t -*-
;;; For more information see (info "(emacs) Directory Variables")

((nil . ((compile-command . "yarn run test")))
 (go-mode
  (dape-configs . ((go-launch-debug
		    modes (go-mode go-ts-mode)
		    ;; ensure dape-ensure-command
		    ;; fn (dape-config-autoport dape-config-tramp)
		    command "dlv"
		    ;; command-insert-stderr t
		    command-args ("dap" "--listen" "127.0.0.1::autoport")
		    command-cwd (lambda () (dape-command-cwd))
		    port :autoport
		    :type "go"
		    :request "launch"
		    :mode "debug"
		    :program "/home/dabanya02/Documents/GoRiichi/cmd/server/main.go"
		    :cwd (lambda () (dape-cwd))
		    :host "127.0.0.1"
		    :showLog t)))))

