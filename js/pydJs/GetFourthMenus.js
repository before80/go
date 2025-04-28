() => {
    // // toctree-wrapper compound
    const tocTreeUl = document.querySelector("div.toctree-wrapper.compound > ul")
    let menuInfos = []
    let exists = {}
    if (tocTreeUl) {
        tocTreeUl.querySelectorAll(":scope > li").forEach(li => {
            const a = li.querySelector(':scope > a')
            const menu_name = a.textContent.trim()
            const url = a.href.trim()
            const names = url.split('/')
            let filename = names[names.length - 1].replace(/\.html$/, '')
                .replace(/[\.\/]/g, '_')
            const noJhaoUrl = url.split("#")[0]
            if (!exists[noJhaoUrl]) {
                menuInfos.push({
                    menu_name: menu_name,
                    filename: filename,
                    url: noJhaoUrl,
                })
                exists[noJhaoUrl] = true
            }
        })
    } else {
        const uls = document.querySelectorAll("ul.simple")
        const curPageUrl = "%s"
        // if (uls.length > 0) {
        //     uls.forEach(ul => {
        //         ul.querySelectorAll(":scope > li").forEach(li => {
        //             const a = li.querySelector('a')
        //             console.log("a=", a)
        //             if (a) {
        //                 const menu_name = a.textContent.trim()
        //                 const url = a.href.trim()
        //                 let names = url.split('/')
        //                 names = names[names.length - 1].split("#")
        //                 let filename = names[0].replace(/\.html$/, '')
        //                     .replace(/[\.\/]/g, '_')
        //                 const noJhaoUrl = url.split("#")[0]
        //                 console.log("noJhaoUrl=", noJhaoUrl)
        //                 if (noJhaoUrl !== curPageUrl && !exists[noJhaoUrl]) {
        //                     menuInfos.push({
        //                         menu_name: menu_name,
        //                         filename: filename,
        //                         url: noJhaoUrl,
        //                     })
        //                     exists[noJhaoUrl] = true
        //                 }
        //             }
        //         })
        //     })
        // }
    }
    console.log(menuInfos)
    return menuInfos
}