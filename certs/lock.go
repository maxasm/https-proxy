package certs

import (
	"os"
	"time"
)


type LockStatus struct {
	GenerateNew bool // should you generate a new lock file and thus new certs
	Ok bool // is everything ok, the certs exist and there is no lock file
}


func check_certs(domain string) bool {
	out_dir := certs_dir+domain+"/"
	// check if the fullchain-cert.pem file exists for this domain
	_, err__check_cert := os.Stat(out_dir+"fullchain-cert.pem")
	if err__check_cert != nil {
		return false
	}

	// check if the pr-key.pem file exists for this domain
	_, err__check_key := os.Stat(out_dir+"pr-key.pem")
	if err__check_key != nil {
		return false
	}

	return true
}

func get_lock_status(domain string) LockStatus {
	out_dir := certs_dir+domain+"/" 

	// check if the lock file exists
	file_stats, err__get_file_stats := os.Stat(out_dir+"lock")
	if err__get_file_stats != nil {
		certs_exists := check_certs(domain)
		return LockStatus{Ok: certs_exists}
	}

	// in this case a lock file exists
	// get the modification time of the file
	mod_time := file_stats.ModTime()

	// compare the files modification time to time.Now() + 1 second
	if mod_time.Add(time.Second*1).Compare(time.Now()) <= 0 {
		// this means that the expiry time is past the current time
		// and this file should have already been deleted, but it hasn't
		// which might imply that there was a problem generating the certs.
		wl.Printf("lock filed expired, generating new certs and a new lock file\n")
		return LockStatus{GenerateNew: true}
	} 

	// in this case, the lock file exists but it hasn't `expired` yet
	// this means that there is another connection to the same domain
	// that is creating the certs. The system has to therefore wait for the 
	// certs to be generated. We wait for about 500ms and call the function
	// again recursively

	time.Sleep(time.Millisecond*500)
	return get_lock_status(domain)
}


// create a lock file if one does not exist
// if a lock file exists, update the ModTime to the current time
func update_lock_file(domain string) error {
	out_dir := certs_dir+domain+"/"
	// check if the lock file exists
	_, err__stat := os.Stat(out_dir+"lock")
	if err__stat != nil {
		// the file does not exist, therfore create file
		_, err__create_file := os.OpenFile(out_dir+"lock",os.O_RDWR | os.O_CREATE, 0666)
		if err__create_file != nil {
			return err__create_file 
		}

		return nil
	}

	// update the access time and modification time to time.Now()
	return os.Chtimes(out_dir+"lock", time.Now(), time.Now())
}
