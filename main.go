package main

import (
	"flag"
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	saveEnabled *bool
)

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: *debug,
		MenuOptions: []*astilectron.MenuItemOptions{
			&astilectron.MenuItemOptions{
				Label: astilectron.PtrStr("Pico-go"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label: astilectron.PtrStr("About"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "about", aboutBody, func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
									return
								}
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending about event failed"))
							}
							return
						},
					},
					{Role: astilectron.MenuItemRoleClose},
				},
			},
			&astilectron.MenuItemOptions{
				Label: astilectron.PtrStr("File"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Accelerator: astilectron.NewAccelerator("Command", "n"),
						Label: astilectron.PtrStr("New"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "new", demoSrc, func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if m != nil {
									if err := json.Unmarshal(m.Payload, &s); err != nil {
										astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
										return
									}	
								}
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending new event failed"))
							}
							return
						},
						Type: astilectron.MenuItemTypeCheckbox,
					},
					{
						Accelerator: astilectron.NewAccelerator("Command", "o"),
						Label: astilectron.PtrStr("Open"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "open", "open this", func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if m != nil {
									if err := json.Unmarshal(m.Payload, &s); err != nil {
										astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
										return
									}	
								}
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending open event failed"))
							}
							return
						},
					},
					{
						Accelerator: astilectron.NewAccelerator("Command", "s"),
						Label: astilectron.PtrStr("Save"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "save", "save this", func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if m != nil {
									if err := json.Unmarshal(m.Payload, &s); err != nil {
										astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
										return
									}	
								}
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending save event failed"))
							}
							return
						},
						Enabled: saveEnabled,
					},
					{
						Accelerator: astilectron.NewAccelerator("Command", "+s"),
						Label: astilectron.PtrStr("Save As..."),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "saveAs", "saveAs this", func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if m != nil {
									if err := json.Unmarshal(m.Payload, &s); err != nil {
										astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
										return
									}	
								}
							}); err != nil {
								astilog.Error(errors.Wrap(err, "sending saveAs event failed"))
							}
							return
						},
					},
				},
			},
			&astilectron.MenuItemOptions{
				Label: astilectron.PtrStr("Edit"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label:       astilectron.PtrStr("Cut"),
						Accelerator: &astilectron.Accelerator{"CmdOrCtrl+X"},
						Role:        astilectron.MenuItemRoleCut,
					},
					{
						Label:       astilectron.PtrStr("Copy"),
						Accelerator: &astilectron.Accelerator{"CmdOrCtrl+C"},
						Role:        astilectron.MenuItemRoleCopy,
					},
					{
						Label:       astilectron.PtrStr("Paste"),
						Accelerator: &astilectron.Accelerator{"CmdOrCtrl+V"},
						Role:        astilectron.MenuItemRolePaste,
					},
				},
			},
		},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			// go func() {
			// 	time.Sleep(5 * time.Second)
			// 	if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
			// 		astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
			// 	}
			// }()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(850),
				Width:           astilectron.PtrInt(850),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
