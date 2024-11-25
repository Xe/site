let xesite = ./types/package.dhall

let Job = xesite.Job

let Salary = xesite.Salary

let Stock = xesite.Stock

let StockKind = xesite.StockKind

let Company = xesite.Company

let Location = xesite.Location

let annual = \(rate : Natural) -> Salary::{ amount = rate }

let hourly = \(rate : Natural) -> Salary::{ amount = rate, per = "hour" }

let annualCAD = \(rate : Natural) -> Salary::{ amount = rate, currency = "CAD" }

let mercerIsland =
      Location::{
      , city = "Mercer Island"
      , stateOrProvince = "WA"
      , country = "USA"
      }

let bellevue = mercerIsland // { city = "Bellevue" }

let mountainView =
      Location::{
      , city = "Mountain View"
      , stateOrProvince = "CA"
      , country = "USA"
      , remote = False
      }

let sf = mountainView // { city = "San Fransisco" }

let montreal =
      Location::{
      , city = "Montreal"
      , stateOrProvince = "QC"
      , country = "CAN"
      , remote = False
      }

let ottawa =
      Location::{ city = "Ottawa", stateOrProvince = "ON", country = "CAN" }

let imvu =
      Company::{
      , name = "IMVU"
      , url = Some "https://imvu.com"
      , tagline =
          "a company whose mission is to help people find and communicate with eachother. Their main product is a 3D avatar-based chat client and its surrounding infrastructure allowing creators to make content for the avatars to wear."
      , location = mountainView // { city = "Redwood City" }
      }

let tailscale =
      Company::{
      , name = "Tailscale"
      , url = Some "https://tailscale.com"
      , tagline =
          "a zero config VPN for building secure networks. Install on any device in minutes. Remote access from any network or physical location."
      , location = ottawa // { city = "Toronto" }
      }

let flyio =
      Company::{
      , name = "Fly.io"
      , url = Some "https://fly.io"
      , tagline =
          "A platform to run code close to users. Deploy the same app to 35 datacentres worldwide in seconds."
      , location = Location::{
        , city = "Chicago"
        , stateOrProvince = "IL"
        , country = "USA"
        , remote = True
        }
      }

let xeserv =
      Company::{
      , name = "Xeserv"
      , url = Some "https://xeserv.us"
      , tagline = "Uber for developer marketing consulting"
      , location = ottawa
      }

in  [ Job::{
      , company = xeserv
      , title = "Chief Executive Officer"
      , startDate = "2024-08-30"
      , salary = annualCAD 150000
      , locations = [ ottawa ]
      , highlights =
        [ "Developer relations"
        , "Content creation"
        , "Technical writing"
        , "Video production"
        , "Generative AI consulting"
        ]
      }
    , Job::{
      , company = flyio
      , title = "Senior Technophilosopher"
      , startDate = "2024-01-07"
      , endDate = Some "2024-08-30"
      , leaveReason = Some "mass layoffs"
      , salary = Salary::{ amount = 260000, per = "year", currency = "CAD" }
      , locations = [ ottawa ]
      , highlights =
        [ "Developer marketing"
        , "Developer relations"
        , "Technical writing"
        , "Team building"
        ]
      }
    , Job::{
      , company = flyio
      , title = "Senior Technophilosopher (contract)"
      , contract = True
      , startDate = "2023-11-06"
      , endDate = Some "2024-01-06"
      , leaveReason = Some "got hired full-time"
      , salary = Salary::{ amount = 15000, per = "month", currency = "USD" }
      , locations = [ ottawa ]
      , highlights = [ "Developer marketing" ]
      }
    , Job::{
      , company = tailscale
      , title = "Archmage of Infrastructure II"
      , startDate = "2023-04-06"
      , endDate = Some "2023-10-10"
      , daysWorked = Some 187
      , leaveReason = Some "position eliminated in re-org"
      , salary = annualCAD 203651
      , locations = [ ottawa ]
      , highlights =
        [ "Founding the developer relations team"
        , "Crafting engaging talks and posts to help people master Tailscale"
        , "Mentorship and fostering growth of my coworkers"
        ]
      }
    , Job::{
      , company = tailscale
      , title = "Archmage of Infrastructure"
      , startDate = "2022-03-01"
      , endDate = Some "2023-04-06"
      , daysWorked = Some 401
      , leaveReason = Some "raise"
      , salary = annualCAD 147150
      , locations = [ ottawa ]
      , highlights =
        [ "The first developer relations person at Tailscale"
        , "Public-facing content writing"
        , "Public speaking"
        , "Developing custom integration solutions and supporting them"
        ]
      }
    , Job::{
      , company = tailscale
      , title = "Software Designer"
      , startDate = "2020-12-14"
      , endDate = Some "2022-03-01"
      , daysWorked = Some 442
      , daysBetween = Some 0
      , salary = annualCAD 135000
      , leaveReason = Some "raise"
      , locations = [ montreal // { remote = True }, ottawa ]
      , highlights =
        [ "Go programming"
        , "SQL integrations"
        , "Public-facing content writing"
        , "Customer support"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "Lightspeed POS"
        , url = Some "https://lightspeedhq.com"
        , tagline =
            "a provider of retail, ecommerce and point-of-sale solutions for small and medium scale businesses."
        , location = montreal
        }
      , title = "Expert principal en fiabilité du site"
      , startDate = "2019-05-06"
      , endDate = Some "2020-11-27"
      , daysWorked = Some 540
      , daysBetween = Some 48
      , salary =
              annualCAD 115000
          //  { stock = Some Stock::{ amount = 7500, liquid = True } }
      , leaveReason = Some "quit"
      , locations = [ montreal ]
      , highlights =
        [ "Migration from cloud to cloud"
        , "Work on the cloud platform initiative"
        , "Crafting reliable infrastructure for clients of customers"
        , "Creation of an internally consistent and extensible command line interface for internal tooling"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "Heroku"
        , url = Some "https://heroku.com"
        , tagline =
            "a cloud Platform-as-a-Service (PaaS) that created the term 'platform as a service'. Heroku currently supports several programming languages that are commonly used on the web. Heroku, one of the first cloud platforms, has been in development since June 2007, when it supported only the Ruby programming language, but now supports Java, Node.js, Scala, Clojure, Python, PHP, and Go."
        , location = sf
        }
      , title = "Senior Software Engineer"
      , startDate = "2017-11-13"
      , endDate = Some "2019-03-08"
      , daysWorked = Some 480
      , daysBetween = Some 0
      , salary = annual 150000
      , leaveReason = Some "quit"
      , locations = [ mountainView, bellevue ]
      , highlights =
        [ "JVM Application Metrics"
        , "Go Runtime Metrics Agent"
        , "Other backend fixes and improvements on Threshold Autoscaling and Threshold Alerting"
        , "Public-facing blogpost writing"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "MBO Partners (Heroku)"
        , tagline = "a staffing agency used to contract me for Heroku."
        , location = Location::{
          , city = "Herndon"
          , stateOrProvince = "VA"
          , country = "USA"
          }
        }
      , title = "Consultant"
      , contract = True
      , startDate = "2017-02-13"
      , endDate = Some "2017-11-13"
      , daysWorked = Some 273
      , daysBetween = Some 83
      , salary = hourly 120
      , leaveReason = Some "hired"
      , locations = [ mountainView ]
      }
    , Job::{
      , company = Company::{
        , name = "Backplane.io"
        , defunct = True
        , location = sf
        }
      , title = "Software Engineer"
      , startDate = "2016-08-24"
      , endDate = Some "2016-11-22"
      , daysWorked = Some 90
      , daysBetween = Some 21
      , salary = annual 105000 // { stock = Some Stock::{ amount = 85000 } }
      , leaveReason = Some "terminated"
      , locations = [ sf ]
      , highlights =
        [ "Performance monitoring of production servers"
        , "Continuous deployment and development in Go"
        , "Learning a lot about HTTP/2 and load balancing"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "Pure Storage"
        , url = Some "https://www.purestorage.com/"
        , tagline =
            "a Mountain View, California-based enterprise data flash storage company founded in 2009. It is traded on the NYSE (PSTG)."
        , location = mountainView
        }
      , title = "Member of Technical Staff"
      , startDate = "2016-04-04"
      , endDate = Some "2016-08-03"
      , daysWorked = Some 121
      , daysBetween = Some 3
      , salary =
              annual 135000
          //  { stock = Some Stock::{
                , amount = 5000
                , liquid = True
                , kind = StockKind.Grant
                }
              }
      , leaveReason = Some "quit"
      , locations = [ mountainView ]
      , highlights = [ "Python 2 code maintenance", "Working with Foone" ]
      }
    , Job::{
      , company = imvu
      , title = "Systems Administrator"
      , startDate = "2016-03-08"
      , endDate = Some "2016-04-01"
      , daysWorked = Some 24
      , daysBetween = Some 1
      , salary = annual 105000
      , leaveReason = Some "quit"
      , locations = [ mountainView // { city = "Redwood City" } ]
      }
    , Job::{
      , company = imvu
      , title = "Site Reliability Engineer"
      , startDate = "2015-03-30"
      , endDate = Some "2016-03-07"
      , daysWorked = Some 343
      , daysBetween = Some 49
      , salary = annual 125000 // { stock = Some Stock::{ amount = 20000 } }
      , leaveReason = Some "demoted"
      , locations = [ mountainView ]
      , highlights =
        [ "Wrote up technical designs"
        , "Implemented technical designs on an over 800 machine cluster"
        , "Continuous learning of a lot of very powerful systems and improving upon them when it is needed"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "VTCSecure"
        , url = Some "https://www.vtcsecure.com/"
        , tagline =
            "a company dedicated to helping with custom and standard audio/video conferencing solutions. They specialize in helping the deaf and blind communicate over today's infrastructure without any trouble on their end."
        , location = Location::{
          , city = "Clearwater"
          , stateOrProvince = "FL"
          , country = "USA"
          }
        }
      , title = "Consultant"
      , contract = True
      , startDate = "2014-10-27"
      , endDate = Some "2015-02-09"
      , daysWorked = Some 105
      , daysBetween = Some 12
      , salary = hourly 90
      , leaveReason = Some "contract not renewed"
      , locations = [ mercerIsland ]
      , highlights =
        [ "Started groundwork for a dynamically scalable infrastructure on a project for helping the blind see things"
        , "Developed a prototype of a new website for VTCSecure"
        , "Education on best practices using Docker and CoreOS"
        , "Learning Freeswitch"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "Appen"
        , url = Some "https://appen.com/"
        , tagline =
            "is a company that uses crowdsourcing to have its customers submit tasks to be done, similar to Amazon's Mechanical Turk."
        , location = mountainView // { city = "San Francisco", remote = True }
        }
      , title = "Consultant"
      , contract = True
      , startDate = "2014-09-17"
      , endDate = Some "2014-10-15"
      , daysWorked = Some 28
      , daysBetween = Some 21
      , salary = hourly 90
      , leaveReason = Some "contract not renewed"
      , locations = [ mercerIsland ]
      , highlights =
        [ "Research and development on scalable Linux deployments on AWS via CoreOS and Docker"
        , "Development of in-house tools to speed instance creation"
        , "Laid groundwork on the creation and use of better tools for managing large clusters of CoreOS and Fleet machines"
        ]
      }
    , Job::{
      , company = Company::{
        , name = "OpDemand"
        , defunct = True
        , tagline =
            "the company behind the open source project Deis, a distributed platform-as-a-service (PaaS) designed from the ground up to emulate Heroku but on privately owned servers."
        , location = Location::{
          , city = "Boulder"
          , stateOrProvince = "CO"
          , country = "USA"
          }
        }
      , title = "Software Engineering Intern"
      , startDate = "2014-07-14"
      , endDate = Some "2014-08-27"
      , daysWorked = Some 44
      , daysBetween = Some 189
      , salary = annual 35000
      , leaveReason = Some "terminated"
      , locations = [ mercerIsland ]
      , highlights =
        [ "Built new base image for Deis components"
        , "Research and development on a new builder component"
        ]
      , hideFromResume = True
      }
    , Job::{
      , company = Company::{
        , name = "Symplicity"
        , tagline =
            "a company that provides students with the tools and connections they need to enhance their employability while preparing to succeed in today's job market."
        , url = Some "https://www.symplicity.com"
        , location = Location::{
          , city = "Arlington"
          , stateOrProvince = "VA"
          , country = "USA"
          , remote = False
          }
        }
      , title = "Junior Systems Administrator"
      , startDate = "2013-11-11"
      , endDate = Some "2014-01-06"
      , daysWorked = Some 56
      , salary = annual 50000
      , leaveReason = Some "terminated"
      , locations =
        [ Location::{
          , city = "Arlington"
          , stateOrProvince = "VA"
          , country = "USA"
          , remote = False
          }
        ]
      , highlights = [ "Python message queue processing" ]
      , hideFromResume = True
      }
    ]
