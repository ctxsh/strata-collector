do we look at packaging the send channels?

the discovery to collector relationship is pretty easy since we are using selectors which gives us a list.  The collector to discovery won't have that luxury.  They'll have a spec list - though we could add labels to discovery?

New discovery added:
* Get a list of collectors we are going to send to.
* Update the collectors by adding the discovery ref to the list.

* Status: active/inactive
* Collectors: connected/unconnected (1/1)


* ? How do we get the channels (id?)

Collector added:
* get a list of discovery services, do any of them reference me?  This means that the refs will need to be put into place without getting the data directly from and I'm not sure that I like the idea of a ref to something that does not exitst.
