# Milkyway :milky_way:
Milkyway is a [TOSCA generator](https://docs.oasis-open.org/tosca/TOSCA/v2.0/csd03/TOSCA-v2.0-csd03.html#_Toc56506778). It will help you to use [OASIS TOSCA](https://www.oasis-open.org/committees/tc_home.php?wg_abbrev=tosca) (Topology and Orchestration Specification for Cloud Applications) standard without knowledge of TOSCA.
You only need to put URL link for Ansible Galaxy Role from Github or path to the directory and get ready CSAR Archive to use it in TOSCA Orchestrator.
# Quick start 
## Build
```
# if you have not got already installed golang -> check "realy quick start" section
go build
go install
```
## Usage
```bigquery
milkyway help
## or milkyway -h, --help
milkyway generate -u <ansible_role_url>

```
### Input
We get github url or path to a folder with role. You need make sure, that this role already worked, we are not magicians and could not correct Role errors.
```
├───files:
|   └──main.yaml
├───templates:
|   └──main.yaml
├───meta:
|   └──main.yaml
├───defaults:
|   └──main.yaml
├───vars:
|   └──main.yaml
├───tasks
|   ├───main.yaml
|   └──setup-Debian.yml
├───handlers
|   └──main.yaml
```
### Output
Program generates CSAR archive with below structure:
```
|───TOSCA.meta
|
├───nodetypes:
|   └──main.yaml
├───artifactypes:
|   └──main.yaml
├───capabilitytypes:
|   └──main.yaml
├───definitions:
|   └──main.yaml
├───defaults:
|   └──main.yaml
```
# Realy quick start :stars:
* [install go](https://golang.org/doc/install)
```
git clone https://github.com/princessfruittt/milkyway.git
cd milkyway
go build
go install
 ```
**contacts:**

:iphone: telegram: princessfruittt

:mailbox: email: princessfruittt@yandex.ru

# License