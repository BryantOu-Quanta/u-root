// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"errors"
	"fmt"
	"strings"
)

const TableTypeBIOSInformation TableType = 0

// Much of this is auto-generated. If adding a new type, see README for instructions.

// BIOSInformation is Defined in DSP0134 7.1.
type BIOSInformation struct {
	Table
	Vendor                                 string                  // 04h
	Version                                string                  // 05h
	StartingAddressSegment                 uint16                  // 06h
	ReleaseDate                            string                  // 08h
	ROMSize                                uint8                   // 09h
	Characteristics                        BIOSCharacteristics     // 0Ah
	CharacteristicsExt1                    BIOSCharacteristicsExt1 // 12h
	CharacteristicsExt2                    BIOSCharacteristicsExt2 // 13h
	SystemBIOSMajorRelease                 uint8                   // 14h
	SystemBIOSMinorRelease                 uint8                   // 15h
	EmbeddedControllerFirmwareMajorRelease uint8                   // 16h
	EmbeddedControllerFirmwareMinorRelease uint8                   // 17h
	ExtendedROMSize                        uint16                  // 18h
}

// BIOSCharacteristics is defined in DSP0134 7.1.1.
type BIOSCharacteristics uint64

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.1.
type BIOSCharacteristicsExt1 uint8

// BIOSCharacteristicsExt2 is defined in DSP0134 7.1.2.2.
type BIOSCharacteristicsExt2 uint8

// NewBIOSInformation parses a generic Table into BIOSInformation.
func NewBIOSInformation(t *Table) (*BIOSInformation, error) {
	if t.Type != TableTypeBIOSInformation {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x12 {
		return nil, errors.New("required fields missing")
	}
	bi := &BIOSInformation{Table: *t}
	if _, err := parseStruct(t, 0 /* off */, false /* complete */, bi); err != nil {
		return nil, err
	}
	return bi, nil
}

// GetROMSizeBytes returns ROM size in bytes.
func (bi *BIOSInformation) GetROMSizeBytes() uint {
	if bi.ROMSize != 0xff {
		return 65536 * (uint(bi.ROMSize) + 1)
	}
	var extSize uint
	if bi.Len() >= 0x1a {
		extSize = uint(bi.ExtendedROMSize)
	} else {
		extSize = 0x10 // 16 MB
	}
	unit := (extSize >> 14)
	multiplier := uint(1)
	switch unit {
	case 0:
		multiplier = 1024 * 1024
	case 1:
		multiplier = 1024 * 1024 * 1024
	}
	return (extSize & 0x3fff) * multiplier
}

func (bi *BIOSInformation) String() string {
	lines := []string{
		bi.Header.String(),
		fmt.Sprintf("\tVendor: %s", bi.Vendor),
		fmt.Sprintf("\tVersion: %s", bi.Version),
		fmt.Sprintf("\tRelease Date: %s", bi.ReleaseDate),
	}
	if bi.StartingAddressSegment != 0 {
		lines = append(lines,
			fmt.Sprintf("\tAddress: 0x%04X0", bi.StartingAddressSegment),
			fmt.Sprintf("\tRuntime Size: %d kB", ((0x10000-int(bi.StartingAddressSegment))<<4)/1024),
		)
	}
	bs := bi.GetROMSizeBytes()
	switch {
	case bs < 1024:
		lines = append(lines, fmt.Sprintf("\tROM Size: %d", bs))
	case bs < 1024*1024:
		lines = append(lines, fmt.Sprintf("\tROM Size: %d kB", bs/1024))
	case bs < 1024*1024*1024:
		lines = append(lines, fmt.Sprintf("\tROM Size: %d MB", bs/1024/1024))
	default:
		lines = append(lines, fmt.Sprintf("\tROM Size: %d GB", bs/1024/1024/1024))
	}
	lines = append(lines,
		fmt.Sprintf("\tCharacteristics:\n%s", bi.Characteristics),
		fmt.Sprintf("%s", bi.CharacteristicsExt1),
		fmt.Sprintf("%s", bi.CharacteristicsExt2),
	)
	if bi.Len() >= 0x16 && bi.SystemBIOSMajorRelease != 0xff { // 2.4+
		lines = append(lines, fmt.Sprintf("\tBIOS Revision: %d.%d", bi.SystemBIOSMajorRelease, bi.SystemBIOSMinorRelease))
	}
	if bi.Len() >= 0x18 && bi.EmbeddedControllerFirmwareMajorRelease != 0xff {
		lines = append(lines, fmt.Sprintf("\tFirmware Revision: %d.%d", bi.EmbeddedControllerFirmwareMajorRelease, bi.EmbeddedControllerFirmwareMinorRelease))
	}
	return strings.Join(lines, "\n")
}

// BIOSCharacteristics fields are defined in DSP0134 7.1.1.
const (
	BIOSCharacteristicsReserved                                                           BIOSCharacteristics = (1 << 0)  // Reserved.
	BIOSCharacteristicsReserved2                                                                              = (1 << 1)  // Reserved.
	BIOSCharacteristicsUnknown                                                                                = (1 << 2)  // Unknown.
	BIOSCharacteristicsBIOSCharacteristicsAreNotSupported                                                     = (1 << 3)  // BIOS Characteristics are not supported.
	BIOSCharacteristicsISAIsSupported                                                                         = (1 << 4)  // ISA is supported.
	BIOSCharacteristicsMCAIsSupported                                                                         = (1 << 5)  // MCA is supported.
	BIOSCharacteristicsEISAIsSupported                                                                        = (1 << 6)  // EISA is supported.
	BIOSCharacteristicsPCIIsSupported                                                                         = (1 << 7)  // PCI is supported.
	BIOSCharacteristicsPCCardPCMCIAIsSupported                                                                = (1 << 8)  // PC card (PCMCIA) is supported.
	BIOSCharacteristicsPlugAndPlayIsSupported                                                                 = (1 << 9)  // Plug and Play is supported.
	BIOSCharacteristicsAPMIsSupported                                                                         = (1 << 10) // APM is supported.
	BIOSCharacteristicsBIOSIsUpgradeableFlash                                                                 = (1 << 11) // BIOS is upgradeable (Flash).
	BIOSCharacteristicsBIOSShadowingIsAllowed                                                                 = (1 << 12) // BIOS shadowing is allowed.
	BIOSCharacteristicsVLVESAIsSupported                                                                      = (1 << 13) // VL-VESA is supported.
	BIOSCharacteristicsESCDSupportIsAvailable                                                                 = (1 << 14) // ESCD support is available.
	BIOSCharacteristicsBootFromCDIsSupported                                                                  = (1 << 15) // Boot from CD is supported.
	BIOSCharacteristicsSelectableBootIsSupported                                                              = (1 << 16) // Selectable boot is supported.
	BIOSCharacteristicsBIOSROMIsSocketed                                                                      = (1 << 17) // BIOS ROM is socketed.
	BIOSCharacteristicsBootFromPCCardPCMCIAIsSupported                                                        = (1 << 18) // Boot from PC card (PCMCIA) is supported.
	BIOSCharacteristicsEDDSpecificationIsSupported                                                            = (1 << 19) // EDD specification is supported.
	BIOSCharacteristicsInt13hJapaneseFloppyForNEC980012MB351KBytessector360RPMIsSupported                     = (1 << 20) // Japanese floppy for NEC 9800 1.2 MB (3.5”, 1K bytes/sector, 360 RPM) is
	BIOSCharacteristicsInt13hJapaneseFloppyForToshiba12MB35360RPMIsSupported                                  = (1 << 21) // Japanese floppy for Toshiba 1.2 MB (3.5”, 360 RPM) is supported.
	BIOSCharacteristicsInt13h525360KBFloppyServicesAreSupported                                               = (1 << 22) // 5.25” / 360 KB floppy services are supported.
	BIOSCharacteristicsInt13h52512MBFloppyServicesAreSupported                                                = (1 << 23) // 5.25” /1.2 MB floppy services are supported.
	BIOSCharacteristicsInt13h35720KBFloppyServicesAreSupported                                                = (1 << 24) // 3.5” / 720 KB floppy services are supported.
	BIOSCharacteristicsInt13h35288MBFloppyServicesAreSupported                                                = (1 << 25) // 3.5” / 2.88 MB floppy services are supported.
	BIOSCharacteristicsInt5hPrintScreenServiceIsSupported                                                     = (1 << 26) // Int 5h, print screen Service is supported.
	BIOSCharacteristicsInt9h8042KeyboardServicesAreSupported                                                  = (1 << 27) // Int 9h, 8042 keyboard services are supported.
	BIOSCharacteristicsInt14hSerialServicesAreSupported                                                       = (1 << 28) // Int 14h, serial services are supported.
	BIOSCharacteristicsInt17hPrinterServicesAreSupported                                                      = (1 << 29) // Int 17h, printer services are supported.
	BIOSCharacteristicsInt10hCGAMonoVideoServicesAreSupported                                                 = (1 << 30) // Int 10h, CGA/Mono Video Services are supported.
	BIOSCharacteristicsNECPC98                                                                                = (1 << 31) // NEC PC-98.
)

func (v BIOSCharacteristics) String() string {
	var lines []string
	if v&BIOSCharacteristicsReserved != 0 {
		lines = append(lines, "\t\tReserved")
	}
	if v&BIOSCharacteristicsReserved2 != 0 {
		lines = append(lines, "\t\tReserved")
	}
	if v&BIOSCharacteristicsUnknown != 0 {
		lines = append(lines, "\t\tUnknown")
	}
	if v&BIOSCharacteristicsBIOSCharacteristicsAreNotSupported != 0 {
		lines = append(lines, "\t\tBIOS characteristics not supported")
	}
	if v&BIOSCharacteristicsISAIsSupported != 0 {
		lines = append(lines, "\t\tISA is supported")
	}
	if v&BIOSCharacteristicsMCAIsSupported != 0 {
		lines = append(lines, "\t\tMCA is supported")
	}
	if v&BIOSCharacteristicsEISAIsSupported != 0 {
		lines = append(lines, "\t\tEISA is supported")
	}
	if v&BIOSCharacteristicsPCIIsSupported != 0 {
		lines = append(lines, "\t\tPCI is supported")
	}
	if v&BIOSCharacteristicsPCCardPCMCIAIsSupported != 0 {
		lines = append(lines, "\t\tPC Card (PCMCIA) is supported")
	}
	if v&BIOSCharacteristicsPlugAndPlayIsSupported != 0 {
		lines = append(lines, "\t\tPNP is supported")
	}
	if v&BIOSCharacteristicsAPMIsSupported != 0 {
		lines = append(lines, "\t\tAPM is supported")
	}
	if v&BIOSCharacteristicsBIOSIsUpgradeableFlash != 0 {
		lines = append(lines, "\t\tBIOS is upgradeable")
	}
	if v&BIOSCharacteristicsBIOSShadowingIsAllowed != 0 {
		lines = append(lines, "\t\tBIOS shadowing is allowed")
	}
	if v&BIOSCharacteristicsVLVESAIsSupported != 0 {
		lines = append(lines, "\t\tVLB is supported")
	}
	if v&BIOSCharacteristicsESCDSupportIsAvailable != 0 {
		lines = append(lines, "\t\tESCD support is available")
	}
	if v&BIOSCharacteristicsBootFromCDIsSupported != 0 {
		lines = append(lines, "\t\tBoot from CD is supported")
	}
	if v&BIOSCharacteristicsSelectableBootIsSupported != 0 {
		lines = append(lines, "\t\tSelectable boot is supported")
	}
	if v&BIOSCharacteristicsBIOSROMIsSocketed != 0 {
		lines = append(lines, "\t\tBIOS ROM is socketed")
	}
	if v&BIOSCharacteristicsBootFromPCCardPCMCIAIsSupported != 0 {
		lines = append(lines, "\t\tBoot from PC Card (PCMCIA) is supported")
	}
	if v&BIOSCharacteristicsEDDSpecificationIsSupported != 0 {
		lines = append(lines, "\t\tEDD is supported")
	}
	if v&BIOSCharacteristicsInt13hJapaneseFloppyForNEC980012MB351KBytessector360RPMIsSupported != 0 {
		lines = append(lines, "\t\tJapanese floppy for NEC 9800 1.2 MB is supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt13hJapaneseFloppyForToshiba12MB35360RPMIsSupported != 0 {
		lines = append(lines, "\t\tJapanese floppy for Toshiba 1.2 MB is supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt13h525360KBFloppyServicesAreSupported != 0 {
		lines = append(lines, "\t\t5.25\"/360 kB floppy services are supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt13h52512MBFloppyServicesAreSupported != 0 {
		lines = append(lines, "\t\t5.25\"/1.2 MB floppy services are supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt13h35720KBFloppyServicesAreSupported != 0 {
		lines = append(lines, "\t\t3.5\"/720 kB floppy services are supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt13h35288MBFloppyServicesAreSupported != 0 {
		lines = append(lines, "\t\t3.5\"/2.88 MB floppy services are supported (int 13h)")
	}
	if v&BIOSCharacteristicsInt5hPrintScreenServiceIsSupported != 0 {
		lines = append(lines, "\t\tPrint screen service is supported (int 5h)")
	}
	if v&BIOSCharacteristicsInt9h8042KeyboardServicesAreSupported != 0 {
		lines = append(lines, "\t\t8042 keyboard services are supported (int 9h)")
	}
	if v&BIOSCharacteristicsInt14hSerialServicesAreSupported != 0 {
		lines = append(lines, "\t\tSerial services are supported (int 14h)")
	}
	if v&BIOSCharacteristicsInt17hPrinterServicesAreSupported != 0 {
		lines = append(lines, "\t\tPrinter services are supported (int 17h)")
	}
	if v&BIOSCharacteristicsInt10hCGAMonoVideoServicesAreSupported != 0 {
		lines = append(lines, "\t\tCGA/mono video services are supported (int 10h)")
	}
	if v&BIOSCharacteristicsNECPC98 != 0 {
		lines = append(lines, "\t\tNEC PC-98")
	}
	return strings.Join(lines, "\n")
}

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.1.
const (
	BIOSCharacteristicsExt1ACPIIsSupported               BIOSCharacteristicsExt1 = (1 << 0) // ACPI is supported.
	BIOSCharacteristicsExt1USBLegacyIsSupported                                  = (1 << 1) // USB Legacy is supported.
	BIOSCharacteristicsExt1AGPIsSupported                                        = (1 << 2) // AGP is supported.
	BIOSCharacteristicsExt1I2OBootIsSupported                                    = (1 << 3) // I2O boot is supported.
	BIOSCharacteristicsExt1LS120SuperDiskBootIsSupported                         = (1 << 4) // LS-120 SuperDisk boot is supported.
	BIOSCharacteristicsExt1ATAPIZIPDriveBootIsSupported                          = (1 << 5) // ATAPI ZIP drive boot is supported.
	BIOSCharacteristicsExt11394BootIsSupported                                   = (1 << 6) // 1394 boot is supported.
	BIOSCharacteristicsExt1SmartBatteryIsSupported                               = (1 << 7) // Smart battery is supported.
)

func (v BIOSCharacteristicsExt1) String() string {
	var lines []string
	if v&BIOSCharacteristicsExt1ACPIIsSupported != 0 {
		lines = append(lines, "\t\tACPI is supported")
	}
	if v&BIOSCharacteristicsExt1USBLegacyIsSupported != 0 {
		lines = append(lines, "\t\tUSB legacy is supported")
	}
	if v&BIOSCharacteristicsExt1AGPIsSupported != 0 {
		lines = append(lines, "\t\tAGP is supported")
	}
	if v&BIOSCharacteristicsExt1I2OBootIsSupported != 0 {
		lines = append(lines, "\t\tI2O boot is supported")
	}
	if v&BIOSCharacteristicsExt1LS120SuperDiskBootIsSupported != 0 {
		lines = append(lines, "\t\tLS-120 boot is supported")
	}
	if v&BIOSCharacteristicsExt1ATAPIZIPDriveBootIsSupported != 0 {
		lines = append(lines, "\t\tATAPI Zip drive boot is supported")
	}
	if v&BIOSCharacteristicsExt11394BootIsSupported != 0 {
		lines = append(lines, "\t\tIEEE 1394 boot is supported")
	}
	if v&BIOSCharacteristicsExt1SmartBatteryIsSupported != 0 {
		lines = append(lines, "\t\tSmart battery is supported")
	}
	return strings.Join(lines, "\n")
}

// BIOSCharacteristicsExt1 is defined in DSP0134 7.1.2.2.
const (
	BIOSCharacteristicsExt2BIOSBootSpecificationIsSupported                  BIOSCharacteristicsExt2 = (1 << 0) // BIOS Boot Specification is supported.
	BIOSCharacteristicsExt2FunctionKeyinitiatedNetworkServiceBootIsSupported                         = (1 << 1) // Function key-initiated network service boot is supported.
	BIOSCharacteristicsExt2TargetedContentDistributionIsSupported                                    = (1 << 2) // Enable targeted content distribution.
	BIOSCharacteristicsExt2UEFISpecificationIsSupported                                              = (1 << 3) // UEFI Specification is supported.
	BIOSCharacteristicsExt2SMBIOSTableDescribesAVirtualMachine                                       = (1 << 4) // SMBIOS table describes a virtual machine. (If this bit is not set, no inference can be made
)

func (v BIOSCharacteristicsExt2) String() string {
	var lines []string
	if v&BIOSCharacteristicsExt2BIOSBootSpecificationIsSupported != 0 {
		lines = append(lines, "\t\tBIOS boot specification is supported")
	}
	if v&BIOSCharacteristicsExt2FunctionKeyinitiatedNetworkServiceBootIsSupported != 0 {
		lines = append(lines, "\t\tFunction key-initiated network boot is supported")
	}
	if v&BIOSCharacteristicsExt2TargetedContentDistributionIsSupported != 0 {
		lines = append(lines, "\t\tTargeted content distribution is supported")
	}
	if v&BIOSCharacteristicsExt2UEFISpecificationIsSupported != 0 {
		lines = append(lines, "\t\tUEFI is supported")
	}
	if v&BIOSCharacteristicsExt2SMBIOSTableDescribesAVirtualMachine != 0 {
		lines = append(lines, "\t\tSystem is a virtual machine")
	}
	return strings.Join(lines, "\n")
}
