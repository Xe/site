let xesite = ../types/package.dhall

let Link = xesite.Link

let Location = xesite.Location

let Job = xesite.Job

let Prelude = ../Prelude.dhall

let xe = ../authors/xe.dhall

let resume = ../resume.dhall

let buzzwords =
      let doer = \(item : Text) -> item

      in  Prelude.Text.concatMapSep ", " Text doer resume.buzzwords

let jobHistory =
      let showDate =
            \(job : Job.Type) ->
              let endDate =
                    merge
                      { Some = \(t : Text) -> t, None = "current" }
                      job.endDate

              in  "${job.startDate} - ${endDate}"

      let showLoc =
            \(l : Location.Type) ->
              "${l.city}, ${l.stateOrProvince}, ${l.country}"

      let workedLocs =
            \(j : Job.Type) ->
              let doer =
                    \(l : Location.Type) -> "\\item Work location: ${showLoc l}"

              in  Prelude.Text.concatMapSep "\n" Location.Type doer j.locations

      let highlights =
            \(j : Job.Type) ->
              let doer = \(t : Text) -> "\\item ${t}"

              in  Prelude.Text.concatMapSep "\n" Text doer j.highlights

      let doer =
            \(job : Job.Type) ->
              ''
              \begin{rSubsection}{${job.company.name}}{${showDate
                                                           job}}{${job.title}}{${showLoc
                                                                                   job.company.location}}

              ${workedLocs job}
              ${highlights job}

              \end{rSubsection}
              ''

      in  Prelude.Text.concatMapSep "\n" Job.Type doer resume.jobs

let publications =
      let doer =
            \(link : Link.Type) ->
              ''
              \begin{rSubsection}{\href{${link.url}}{${link.title}}}{}{}{}

              \item ${link.description}

              \end{rSubsection}
              ''

      in  Prelude.Text.concatMapSep
            "\n"
            Link.Type
            doer
            resume.notablePublications

in  ''
    \documentclass{resume}

    \usepackage[left=0.75in,top=0.6in,right=0.75in,bottom=0.6in]{geometry} % Document margins
    \usepackage{hyperref}
    \newcommand{\tab}[1]{\hspace{.2667\textwidth}\rlap{#1}}
    \newcommand{\itab}[1]{\hspace{0em}\rlap{#1}}
    \name{${xe.name}}
    \address{https://xeiaso.net \\ me@xeiaso.net}

    \begin{document}

    \begin{rSection}{Technical Strengths}

    ${buzzwords}

    \end{rSection}

    \begin{rSection}{Experience}

    ${jobHistory}

    \end{rSection}

    \begin{rSection}{Notable Publications}

    ${publications}

    \end{rSection}

    \end{document}
    ''
