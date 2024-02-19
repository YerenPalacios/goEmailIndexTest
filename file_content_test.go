package main

import (
	"io/ioutil"
	"testing"
)

func AssertEqual(t *testing.T, expected string, current string) {
	if current != expected {
		t.Fatalf("Expected \"%v\", got \"%v\"", current, expected)
	}
}

func GetFileContent(t *testing.T, path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading file %v, %v", path, err)
	}
	return string(file)
}

func TestGetFileContent(t *testing.T) {
	fileMap := getMapContent(GetFileContent(t, "./test_files/3"))
	AssertEqual(t, fileMap["Message-ID"], "<27430435.1075840339331.JavaMail.evans@thyme>")
	AssertEqual(t, fileMap["_id"], "<27430435.1075840339331.JavaMail.evans@thyme>")
	AssertEqual(t, fileMap["Date"], "2001-06-19T12:08:21-07:00")
	AssertEqual(t, fileMap["From"], "liz.legros@enron.com")
	AssertEqual(t, fileMap["To"], "don.baughman@enron.com, test@newline.email, test2@newline.email")
	AssertEqual(t, fileMap["Subject"], "FW: Online Timesheets")
	AssertEqual(t, fileMap["Mime-Version"], "1.0")
	AssertEqual(t, fileMap["Content-Type"], "text/plain; charset=us-ascii")
	AssertEqual(t, fileMap["Content-Transfer-Encoding"], "7bit")
	AssertEqual(t, fileMap["X-From"], "Legros, Liz </O=ENRON/OU=NA/CN=RECIPIENTS/CN=LLEGROS>")
	AssertEqual(t, fileMap["X-To"], "Baughman Jr., Don </O=ENRON/OU=NA/CN=RECIPIENTS/CN=Dbaughm>")
	AssertEqual(t, fileMap["X-cc"], "")
	AssertEqual(t, fileMap["X-bcc"], "")
	AssertEqual(t, fileMap["X-Folder"], `\ExMerge - Baughman Jr., Don\Enron Power\GRI Timesheets`)
	AssertEqual(t, fileMap["X-Origin"], "BAUGHMAN-D")
	AssertEqual(t, fileMap["X-FileName"], `don baughman 6-25-02.PST`)

	expectedContent := ` 
Your department has just placed an order for an admin/clerical temporary worker.  This temporary worker will submit his/her weekly time online for your electronic approval.  Your response to each of these e-mails is necessary for billing and payment of the temporary worker.
 
It is easy to approve time online.  The steps are briefly outlined here.  After you receive the e-mail requesting your approval:
1.      Click on the link in the e-mail to be taken to the Submitted Timesheet page.
2.      Logon to the Submitted Timesheet page to view the timesheet. (logon and password was provided in an earlier e-mail)
3.      Click on 'Approve' or 'Reject' after reviewing the timesheet.  If 'rejected,' you will be required to enter comments.
4.      E-mail will be sent to the temporary worker and to CSS advising them of your approval or rejection.  
 
All timesheets must be approved by close of business each Monday.  It is strongly recommended that you review electronic timesheets as soon as they are submitted to allow for resolution of any issues associated with the timesheets prior to the deadline.  Any rejected timesheets will be revised by the Temporary Worker and resubmitted for approval.
 
Thank you for your support of this transition.  We invite you to direct any questions or comments regarding the online timekeeping procedure for admin/clerical temporary workers to Gwen Chavis / Craig McGee /Liz LeGros at 713-345-6899.
 
Liz LeGros
Contigent Resource Consultant
CSS @ Enron
713/345-6899
Fax:  713/646-6186
liz.legros@enron.com
 `
	AssertEqual(t, fileMap["content"], expectedContent)
}
