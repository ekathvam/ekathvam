package nginx

import (
	"fmt"
	"os"
	"ekathvam/shared"
	"ekathvam/shared"
)

func Nginx () {
	pm := globalshared.Get_pm()
	fmt.Println("Package manager is found as ", pm)
	
	if pm == "unknown" {
		fmt.Println("Unable to find the Package Manager")
		os.Exit(1)
	}
	fmt.Println("Installing or Updating the nginx.")

	err := globalshared.Install_pkg(pm, "nginx")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	err = globalshared.Create_dir("/ekathvam/nginx/")
	if err != nil {
		fmt.Println("Failed to create essential directories..")
		os.Exit(1)
	}

}




