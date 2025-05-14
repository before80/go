() => {
    // const btn = document.querySelector("button.js-showInternalDirectories")
    // if (btn) {
    //     if (btn.textContent.trim() === "Show internal") {
    //         btn.click()
    //         setTimeout(function (){}, 1000)
    //     }
    // }
    const table = document.querySelector("table.UnitDirectories-table")
    let pkgInfos = []
    let curTopFilename = ""
    if (!table) {
        return
    }
    const trs = table.querySelectorAll(":scope > tbody > tr")
    if (trs.length === 0) {
        return
    }
    let topMenuIndex = 0
    let subMenuIndex = 0
    trs.forEach(tr => {
        console.log("run here 2")
        let desc = ""
        let menuName = ""
        let filename = ""
        let url = ""
        let isTopMenu = 2
        if (!tr.classList.contains("UnitDirectories-internal")) {
            const td1 = tr.querySelector(":scope > td:first-child")

            if (td1) {
                const div1 = td1.querySelector(":scope > div:first-child")
                if (div1) {
                    if (div1.getAttribute("class") === "UnitDirectories-subdirectory") {
                        if (td1.querySelector("a")) {
                            menuName = td1.querySelector("a").textContent.trim()
                            filename = menuName.replace(/[\/\s]/g, '_')
                            url = td1.querySelector("a").href.trim()
                            subMenuIndex++
                        }
                    } else {
                        if (td1.querySelector("a")) {
                            menuName = td1.querySelector("a").textContent.trim()
                            filename = menuName.replace(/[\/\s]/g, '_')
                            url = td1.querySelector("a").href.trim()
                            curTopFilename = filename
                            isTopMenu = 1
                        } else {
                            if (td1.querySelector("span")) {
                                menuName = td1.querySelector("span").textContent.trim()
                                filename = menuName.replace(/[\/\s]/g, '_')
                                url = ""
                                curTopFilename = filename
                                isTopMenu = 1
                            }
                        }
                    }
                }
            }
            const td2 = tr.querySelector(":scope > td:nth-child(2)")
            if (td2) {
                desc = td2.textContent.trim()
            }
            if (filename && !(filename.startsWith("internal") || filename.indexOf("internal") != -1)) {
                if (isTopMenu === 1) {
                    subMenuIndex = 0
                    if (topMenuIndex === 0) {
                        topMenuIndex++
                    }
                    pkgInfos.push({
                        menu_name: menuName,
                        filename: filename,
                        url: url,
                        desc: desc,
                        is_top_menu: isTopMenu,
                        weight: topMenuIndex,
                        p_filename: "",
                        children: [],
                    })
                } else {
                    pkgInfos.push({
                        menu_name: menuName,
                        filename: filename,
                        url: url,
                        desc: desc,
                        is_top_menu: isTopMenu,
                        weight: subMenuIndex,
                        p_filename: curTopFilename,
                        children: [],
                    })
                    for (let i = 0; i < pkgInfos.length; i++) {
                        if (pkgInfos[i].filename === curTopFilename) {
                            pkgInfos[i].children.push(filename);
                        }
                    }
                }
            }
        }
    })

    console.log("pkgInfos=", pkgInfos)
    return pkgInfos
}