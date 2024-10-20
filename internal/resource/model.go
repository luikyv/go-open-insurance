package resource

import "github.com/luikyv/go-open-insurance/internal/api"

var resourceTypes = []api.ResourceType{
	api.ResourceTypeCAPITALIZATIONTITLES,
	api.ResourceTypeCAPITALIZATIONTITLEWITHDRAWAL,
	api.ResourceTypeCLAIMNOTIFICATION,
	api.ResourceTypeCONTRACTLIFEPENSION,
	api.ResourceTypeCONTRACTPENSIONPLAN,
	api.ResourceTypeCUSTOMERSBUSINESSADDITIONALINFO,
	api.ResourceTypeCUSTOMERSBUSINESSIDENTIFICATIONS,
	api.ResourceTypeCUSTOMERSBUSINESSQUALIFICATION,
	api.ResourceTypeCUSTOMERSPERSONALADDITIONALINFO,
	api.ResourceTypeCUSTOMERSPERSONALIDENTIFICATIONS,
	api.ResourceTypeCUSTOMERSPERSONALQUALIFICATION,
	api.ResourceTypeDAMAGESANDPEOPLEACCEPTANCEANDBRANCHESABROAD,
	api.ResourceTypeDAMAGESANDPEOPLEAUTO,
	api.ResourceTypeDAMAGESANDPEOPLEFINANCIALRISKS,
	api.ResourceTypeDAMAGESANDPEOPLEHOUSING,
	api.ResourceTypeDAMAGESANDPEOPLEPATRIMONIAL,
	api.ResourceTypeDAMAGESANDPEOPLEPERSON,
	api.ResourceTypeDAMAGESANDPEOPLERESPONSIBILITY,
	api.ResourceTypeDAMAGESANDPEOPLERURAL,
	api.ResourceTypeDAMAGESANDPEOPLETRANSPORT,
	api.ResourceTypeENDORSEMENT,
	api.ResourceTypeFINANCIALASSISTANCE,
	api.ResourceTypeLIFEPENSION,
	api.ResourceTypePENSIONPLAN,
	api.ResourceTypePENSIONWITHDRAWAL,
	api.ResourceTypeQUOTEACCEPTANCEANDBRANCHESABROAD,
	api.ResourceTypeQUOTEAUTO,
	api.ResourceTypeQUOTECAPITALIZATIONTITLE,
	api.ResourceTypeQUOTEFINANCIALRISK,
	api.ResourceTypeQUOTEHOUSING,
	api.ResourceTypeQUOTEPATRIMONIAL,
	api.ResourceTypeQUOTEPERSON,
	api.ResourceTypeQUOTERESPONSIBILITY,
	api.ResourceTypeQUOTERURAL,
	api.ResourceTypeQUOTETRANSPORT,
}