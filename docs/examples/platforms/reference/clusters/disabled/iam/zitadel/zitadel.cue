package holos

#TargetNamespace: #InstancePrefix + "-zitadel"

#DB: {
	Host: "crdb-public"
}

// The canonical login domain for the entire platform.  Zitadel will be active on a singlec cluster at a time, but always accessible from this hostname.
#ExternalDomain: "login.\(#Platform.org.domain)"
