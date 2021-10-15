package model

const AppleFilesDefaultPath string = "./apple"

// AppleFiles holds together static files needed for supporting Apple services.
type AppleFiles struct {
	DeveloperDomainAssociation string `yaml:"developerDomainAssociation" json:"developer_domain_association"`
	AppSiteAssociation         string `yaml:"appSiteAssociation" json:"app_site_association"`
}

// AppleFilenames are names of the files related to Apple services.
var AppleFilenames = AppleFiles{
	DeveloperDomainAssociation: "apple-developer-domain-association.txt",
	AppSiteAssociation:         "apple-app-site-association",
}
