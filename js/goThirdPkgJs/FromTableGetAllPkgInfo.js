function FromTableGetAllPkgInfo(topsidePkgName, weight) {
    const table = document.querySelector("table.UnitDirectories-table")
    let pkgInfos = []
    pkgInfos.push({
        pkg_name: topsidePkgName,
        filename: "_index",
        url: window.location.href,
        dir: topsidePkgName,
        weight: weight,
        desc: "",
    })
    let pkgName = ""
    let totalPkgName = ""
    let url = ""
    let dir = ""
    let desc = ""
    let curWeight = weight
    let ignoreTds = []
    let curTopPkgName = ""
    let filenameUseIndex = false
    if (table) {
        const trs = table.querySelectorAll(":scope > tbody > tr")
        if (trs.length > 0) {
            trs.forEach(tr => {
                filenameUseIndex = false
                const td1 = tr.querySelector(":scope > td:first-child")
                const td2 = tr.querySelector(":scope > td:nth-child(2)")
                if (td1) {
                    const btn = td1.querySelector("button")
                    if (btn) {
                        btn.remove()
                    }
                    const mobileSynopsis = td1.querySelector(".UnitDirectories-mobileSynopsis")
                    if (mobileSynopsis) {
                        mobileSynopsis.remove()
                    }
                    totalPkgName = td1.textContent.trim()
                    if (totalPkgName.startsWith("_")) {
                        const ariaCStr = tr.getAttribute("aria-controls")
                        const ariaCTemps = ariaCStr.split(" ")
                        // console.log("ariaCTemps=",ariaCTemps)
                        for (let aria of ariaCTemps) {
                            ignoreTds.push(aria.replace(`${totalPkgName}-`, ""))
                        }
                    } else {
                        if (td1.querySelector("div.UnitDirectories-pathCell")) {
                            filenameUseIndex = true
                            curTopPkgName = totalPkgName
                            dir = `${topsidePkgName}`
                        }


                        if (td1.querySelector("a")) {
                            url = td1.querySelector("a").href.trim()
                            pkgName = td1.querySelector("a").textContent.trim()
                            if (pkgName.includes("/")) {
                                let names = pkgName.split("/")
                                dir = `${topsidePkgName}/${curTopPkgName}`
                                for (let i = 0; i <= names.length - 2; i++) {
                                    dir =  `${dir}/${names[i]}`
                                }
                                pkgName = names[names.length - 1]
                            }

                            if (td2) {
                                desc = td2.textContent.trim()
                            } else {
                                desc = ""
                            }

                            if (!ignoreTds.includes(totalPkgName)) {
                                curWeight = curWeight + 1
                                pkgInfos.push({
                                    pkg_name: pkgName,
                                    filename: filenameUseIndex ? "_index" : pkgName,
                                    url: url,
                                    dir: dir,
                                    weight: curWeight,
                                    desc: desc,
                                })
                            }
                        }
                    }
                }
            })
        }
    }

    console.log(pkgInfos)
    for (let obj of pkgInfos) {
        console.log(`${obj.url}||${obj.pkg_name}||${obj.dir}||${obj.filename}||${obj.weight}||${obj.desc}`)
    }

    // return pkgInfos
}

FromTableGetAllPkgInfo("chi", 5020)