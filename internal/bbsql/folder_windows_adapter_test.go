//go:build windows

package bbsql

import (
	"testing"

	"golang.org/x/sys/windows"
)

// applySDDL sets an explicit DACL on a folder the test owns, the same way
// install-sql does. Nothing outside the test's own temp folder is touched.
func applySDDL(t *testing.T, dir, sddl string) {
	t.Helper()
	sd, err := windows.SecurityDescriptorFromString(sddl)
	if err != nil {
		t.Fatalf("parse %q: %v", sddl, err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("dacl of %q: %v", sddl, err)
	}
	// Honour the SDDL's own protection flag: forcing PROTECTED here would make an
	// intentionally unprotected descriptor impossible to test.
	ctrl, _, err := sd.Control()
	if err != nil {
		t.Fatalf("control of %q: %v", sddl, err)
	}
	var info windows.SECURITY_INFORMATION = windows.DACL_SECURITY_INFORMATION
	if ctrl&windows.SE_DACL_PROTECTED != 0 {
		info |= windows.PROTECTED_DACL_SECURITY_INFORMATION
	}
	if err := windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, info,
		nil, nil, dacl, nil); err != nil {
		t.Fatalf("apply %q: %v", sddl, err)
	}
}

// W10-01 regression: the installer's own explicitly locked folder must be
// ACCEPTED. x/sys DACL() returns "defaulted", not "present"; reading that as
// presence made every correct install look unlocked and refused startup. The
// pre-existing smoke test only asserted refusal of a temp folder, so it passed
// even while the reader always refused.
func TestFolderAdapterAcceptsInstallerLock(t *testing.T) {
	self := folderSelfSID()
	if self == "" {
		t.Skip("this process token has no readable SID")
	}
	dir := t.TempDir()
	// SYSTEM + Administrators + this account, protected, nothing inherited: the
	// shape install-sql writes, with the running account as the service trustee.
	applySDDL(t, dir, "D:PAI(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)(A;OICI;FA;;;"+self+")")

	if err := CheckFolderAccess(dir); err != nil {
		t.Fatalf("the installer's locked folder was refused: %v", err)
	}
	acl, err := readFolderACL(dir)
	if err != nil {
		t.Fatalf("readFolderACL: %v", err)
	}
	if !acl.Present || !acl.Protected || acl.Unknown || acl.Inherited {
		t.Fatalf("adapter read %+v, want present+protected, nothing unknown or inherited", acl)
	}
	if len(acl.Grantees) != 3 {
		t.Fatalf("grantees = %v, want the three trustees", acl.Grantees)
	}
}

// W10-12: an allow form whose SID is not at SidStart must refuse, never be
// mistaken for a deny and silently ignored. Type 11 is
// ACCESS_ALLOWED_CALLBACK_OBJECT_ACE, not a deny.
func TestFolderAdapterRefusesUnreadableGrantForms(t *testing.T) {
	self := folderSelfSID()
	if self == "" {
		t.Skip("this process token has no readable SID")
	}
	base := "D:PAI(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)(A;OICI;FA;;;" + self + ")"
	for _, tc := range []struct {
		name string
		sddl string
	}{
		// Conditional object allow (type 11): a grant this check cannot read.
		{"callback object allow", base + "(ZA;;FR;00112233-4455-6677-8899-aabbccddeeff;;WD;(Member_of {SID(S-1-5-32-544)}))"},
		// A plain grant to Everyone is outside the trustee set.
		{"everyone grant", base + "(A;OICI;FA;;;WD)"},
		// An empty DACL grants nobody: the service could not write either.
		{"empty dacl", "D:PAI"},
	} {
		dir := t.TempDir()
		applySDDL(t, dir, tc.sddl)
		if err := CheckFolderAccess(dir); err == nil {
			t.Errorf("%s: accepted, want refusal", tc.name)
		}
	}
}

// An unprotected DACL is refused even when its trustees are exactly right.
func TestFolderAdapterRefusesUnprotected(t *testing.T) {
	self := folderSelfSID()
	if self == "" {
		t.Skip("this process token has no readable SID")
	}
	dir := t.TempDir()
	applySDDL(t, dir, "D:AI(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)(A;OICI;FA;;;"+self+")")
	if err := CheckFolderAccess(dir); err == nil {
		t.Fatal("an unprotected DACL was accepted")
	}
}
