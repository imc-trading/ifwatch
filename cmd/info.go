package cmd

import (
	"log"
	"strings"

	"github.com/imc-trading/ifwatch/netx"

	"github.com/mickep76/go-sff"
	"github.com/mickep76/termtable"
	"github.com/mickep76/termtable/winsize"
)

func Info(args map[string]interface{}) {
	all, err := netx.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	interfaces := netx.InterfaceList{}
	if args["<interface>"] != nil {
		for _, i := range all {
			if i.Name == args["<interface>"].(string) {
				interfaces = append(interfaces, i)
				break
			}
		}

		if len(interfaces) == 0 {
			log.Fatalf("interface doesn't exist: %s", args["<interface>"].(string))
		}
	} else {
		interfaces = all
	}

	ws, _ := winsize.Get()
	t := table.New().SetWidth(ws.Cols)

	c1 := t.AddCol("Name").SetWidthPerc(30)
	c1.Cyan()
	c1.HeadingColor().White()

	c2 := t.AddCol("Driver").SetWidthPerc(10)
	c2.Green()
	c2.HeadingColor().White()

	c3 := t.AddCol("Slot").SetWidthPerc(10)
	c3.Green()
	c3.HeadingColor().White()

	c4 := t.AddCol("Flags").SetWidthPerc(30)
	c4.Yellow()
	c4.HeadingColor().White()

	c5 := t.AddCol("Hw Addr").SetWidthPerc(10)
	c5.Green()
	c5.HeadingColor().White()

	c6 := t.AddCol("IPv4").SetWidthPerc(10)
	c6.Green()
	c6.HeadingColor().White()

	c7 := t.AddCol("Netmask").SetWidthPerc(10)
	c7.Green()
	c7.HeadingColor().White()

	c8 := t.AddCol("Network").SetWidthPerc(10)
	c8.Green()
	c8.HeadingColor().White()

	c9 := t.AddCol("Module").SetWidthPerc(10)
	c9.Green()
	c9.HeadingColor().White()

	c10 := t.AddCol("Vendor").SetWidthPerc(10)
	c10.Green()
	c10.HeadingColor().White()

	c11 := t.AddCol("PN").SetWidthPerc(10)
	c11.Green()
	c11.HeadingColor().White()

	c12 := t.AddCol("Rev").SetWidthPerc(10)
	c12.Green()
	c12.HeadingColor().White()

	c13 := t.AddCol("SN").SetWidthPerc(10)
	c13.Green()
	c13.HeadingColor().White()

	for _, i := range interfaces {
		var module, vendor, vendorPn, vendorRev, vendorSn string
		if i.Module != nil {
			switch i.Module.Type {
			case sff.TypeSff8079:
				module = i.Module.Sff8079.Identifier.String()
				vendor = i.Module.Sff8079.Vendor.String()
				vendorPn = i.Module.Sff8079.VendorPn.String()
				vendorRev = i.Module.Sff8079.VendorRev.String()
				vendorSn = i.Module.Sff8079.VendorSn.String()
			case sff.TypeSff8636:
				module = i.Module.Sff8636.Identifier.String()
				vendor = i.Module.Sff8636.Vendor.String()
				vendorPn = i.Module.Sff8636.VendorPn.String()
				vendorRev = i.Module.Sff8636.VendorRev.String()
				vendorSn = i.Module.Sff8636.VendorSn.String()
			}
		}

		t.AddRow(i.Name, i.Driver, i.Slot, strings.Join(i.Flags, ", "), i.HwAddr, i.IPv4, i.Netmask, i.Network, module, vendor, vendorPn, vendorRev, vendorSn)
	}
	t.SortDesc("Name").Sort().Format().Print()
}
