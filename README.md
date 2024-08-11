# metallist
This is a conceptual description of the software to solve my problems with media server automation that I encounter when using Sonarr and other applications.

Perhaps some of my problems have already been solved by software that I am not aware of.

# Basic outline

Basic flow of the automated media server:
1. Add to list
2. Download files
3. Organize files
4. Serve files

There are a huge number of services and databases, and the problem is that different software allows you to use different sets of services and databases and each in its own way.

## Lists

Services to keep lists and track your media consumption:
* Not all services are supported by the software.
* Progress tracking usually is not automated.
* It is difficult or generally impossible to move your lists when there is a need to switch from one service to another.

### Solution like

* Prowlarr but for lists.

## Metadata

Different databases, different IDs, different structure, different information is requested by different programs from several sources.

### Solution like

* Hub for metadata

## Contributing

It is related to metadata. I haven't figured it out yet, but it would be nice to have a short list of actions for making edits. Moreover, help to make connections between different services/entities.

### Solution like

* Shoko with anidb.com ...
* Contribute to TheXEM
