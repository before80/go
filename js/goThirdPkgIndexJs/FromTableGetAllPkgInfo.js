() => {
    let topsidePkgName = "%s"
    let weight = "%d"
    const table = document.querySelector("table.UnitDirectories-table")
    let pkgInfos = []
    pkgInfos.push({
        pkg_name: topsidePkgName,
        filename: "_index",
        url: window.location.href,
        dir: topsidePkgName,
        need_pre_create_index: 2,
        pre_create_index_weight: 0,
        weight: parseInt(weight),
        desc: "",
    })
    let pkgName = ""
    let completePkgName = ""
    let url = ""
    let dir = ""
    let desc = ""
    let curTopPkgWeight = 0
    let curSubPkgWeight = 0
    let ignoreTds = []
    let includeTds = []
    let curTopPkgName = ""
    let filenameUseIndex = false
    let needPreCreateIndexTrNum = -1;
    let needPreCreateIndex = false
    let hadBtn = false
    let hadSubHadA = false
    let resetSubPkgWeight = false
    let preCreateIndexWeight = 0
    if (table) {
        const trs = table.querySelectorAll(":scope > tbody > tr")
        if (trs.length > 0) {
            trs.forEach((tr,trNum) => {
                if (!tr.classList.contains("UnitDirectories-internal")) {
                    dir = `${topsidePkgName}`
                    filenameUseIndex = false
                    hadBtn = false
                    hadSubHadA = false
                    const td1 = tr.querySelector(":scope > td:first-child")
                    const td2 = tr.querySelector(":scope > td:nth-child(2)")
                    if (td1) {
                        const btn = td1.querySelector("button")
                        if (btn) {
                            hadBtn = true
                            btn.remove()
                        }
                        const mobileSynopsis = td1.querySelector(".UnitDirectories-mobileSynopsis")
                        if (mobileSynopsis) {
                            mobileSynopsis.remove()
                        }
                        completePkgName = td1.textContent.trim()
                        if (completePkgName.startsWith("_") || completePkgName === "internal" || completePkgName === "examples") {
                            // const ariaCStr = tr.getAttribute("aria-controls")
                            // if (ariaCStr) {
                            //     const ariaCTemps = ariaCStr.trim().split(" ")
                            //     // console.log("ariaCTemps=",ariaCTemps)
                            //     for (let aria of ariaCTemps) {
                            //         ignoreTds.push(aria.replace(`${totalPkgName}-`, ""))
                            //     }
                            // }
                            const ariaOwnStr = td1.getAttribute("data-aria-owns")
                            if (ariaOwnStr) {
                                const ariaOwnTemps = ariaOwnStr.trim().split(" ")
                                // console.log("ariaCTemps=",ariaCTemps)
                                for (let aria of ariaOwnTemps) {
                                    ignoreTds.push(aria.replace(`${completePkgName}-`, ""))
                                }
                            }
                        } else {
                            if (td1.querySelector("div.UnitDirectories-pathCell")) {
                                if (hadBtn) {
                                    filenameUseIndex = true
                                    curTopPkgName = completePkgName
                                    dir = `${topsidePkgName}/${completePkgName}`
                                } else {
                                    filenameUseIndex = false
                                }
                                const ariaOwnStr = td1.getAttribute("data-aria-owns")
                                if (ariaOwnStr != "") {
                                    includeTds = []
                                    const ariaOwnTemps = ariaOwnStr.trim().split(" ")
                                    // console.log("ariaCTemps=",ariaCTemps)
                                    for (let aria of ariaOwnTemps) {
                                        includeTds.push(aria.replace(`${completePkgName}-`, ""))
                                    }
                                    resetSubPkgWeight = true
                                    curTopPkgWeight = curTopPkgWeight + 1
                                    preCreateIndexWeight = curTopPkgWeight
                                } else {
                                    resetSubPkgWeight = false
                                    preCreateIndexWeight = 0
                                }

                                if (!td1.querySelector("a")) {
                                    needPreCreateIndexTrNum = trNum+1
                                } else {
                                    hadSubHadA = true
                                }
                            }


                            if (td1.querySelector("a")) {
                                url = td1.querySelector("a").href.trim()
                                pkgName = td1.querySelector("a").textContent.trim()
                                let needAddCurTopPkgName = false

                                if (needPreCreateIndexTrNum === trNum) {
                                    needPreCreateIndex = true
                                } else {
                                    needPreCreateIndex = false
                                }

                                if (includeTds.includes(pkgName) || needPreCreateIndexTrNum === trNum || hadSubHadA) {
                                    needAddCurTopPkgName = true
                                }
                                console.log("includeTds=",includeTds,"pkgName=", pkgName,"needAddCurTopPkgName=", needAddCurTopPkgName)
                                if (needAddCurTopPkgName && curTopPkgName !== "") {
                                    dir = `${topsidePkgName}/${curTopPkgName}`
                                } else {
                                    dir = `${topsidePkgName}`
                                }
                                if (pkgName.includes("/")) {
                                    let names = pkgName.split("/")

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

                                if (!ignoreTds.includes(completePkgName) && completePkgName.indexOf("internal") == -1 && completePkgName.indexOf("examples") == -1) {
                                    if (resetSubPkgWeight) {
                                        curSubPkgWeight = 0
                                        curSubPkgWeight = curSubPkgWeight + 1
                                        pkgInfos.push({
                                            pkg_name: pkgName,
                                            filename: filenameUseIndex ? "_index" : pkgName,
                                            url: url,
                                            dir: dir,
                                            weight: curSubPkgWeight,
                                            need_pre_create_index: trNum === needPreCreateIndexTrNum || needPreCreateIndex ? 1 : 2,
                                            pre_create_index_weight: preCreateIndexWeight,
                                            desc: desc,
                                        })
                                    } else {
                                        curTopPkgWeight = curTopPkgWeight + 1
                                        pkgInfos.push({
                                            pkg_name: pkgName,
                                            filename: filenameUseIndex ? "_index" : pkgName,
                                            url: url,
                                            dir: dir,
                                            weight: curTopPkgWeight,
                                            need_pre_create_index: trNum === needPreCreateIndexTrNum || needPreCreateIndex ? 1 : 2,
                                            pre_create_index_weight: preCreateIndexWeight,
                                            desc: desc,
                                        })
                                    }
                                }
                            }
                        }
                    }
                }
            })
        }
    }

    console.log(pkgInfos)
    return pkgInfos
}
