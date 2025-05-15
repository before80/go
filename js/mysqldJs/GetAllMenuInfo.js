() => {
    let menuInfos = []
    function getNestedLiData(ulElement) {
        const liData = [];
        const lis = ulElement.children;
        for (let i = 0; i < lis.length; i++) {
            const li = lis[i];
            if (li.tagName === 'LI') {
                const a = li.querySelector(":scope div.docs-sidebar-nav-link > a")
                if (a) {
                    const url = a.href.trim()
                    const urls = url.split("/")
                    const filename = urls[urls.length - 1].replace(/\.html$/, "")
                    const data = {
                        menu_name: a.textContent.trim().replace(/[\"\'\/\\\$#@&\(\)]/g,'').replace(/\s+/g, " "),
                        filename: filename,
                        url: url,
                        weight: i + 1,
                        children: [],
                    };
                    const nestedUl = li.querySelector(':scope > div.docs-submenu > ul');
                    if (nestedUl) {
                        data.children = getNestedLiData(nestedUl)
                    }
                    liData.push(data);
                }
            }
        }
        return liData;
    }

    function pushChildrenIData(chd, parentDir) {
        const iData = chd
        for (let i = 0; i < iData.length; i++) {
            const dv = iData[i]
            if (dv.children && dv.children.length > 0) {
                menuInfos.push({
                    menu_name: dv.menu_name,
                    is_top_menu: 2,
                    filename: dv.filename,
                    url: dv.url,
                    dir: `${parentDir}/${dv.filename}`,
                    have_sub: 1,
                    weight: dv.weight,
                })
                pushChildrenIData(dv.children, `${parentDir}/${dv.filename}`)
            } else {
                menuInfos.push({
                    menu_name: dv.menu_name,
                    is_top_menu: 2,
                    filename: dv.filename,
                    url: dv.url,
                    dir: parentDir,
                    have_sub: 2,
                    weight: dv.weight,
                })
            }
        }
    }

    function GetAllMenuInfo() {
        let tempMenuInfos = []
        tempMenuInfos = getNestedLiData(document.querySelector("#doc-201 > ul"))
        console.log("tempMenuInfos=",tempMenuInfos)
        if (tempMenuInfos.length > 0) {
            for (let i = 0; i < tempMenuInfos.length; i++) {
                const dv = tempMenuInfos[i]
                if (dv.children && dv.children.length > 0) {
                    menuInfos.push({
                        menu_name: dv.menu_name,
                        is_top_menu: 1,
                        filename: dv.filename,
                        url: dv.url,
                        dir: dv.filename,
                        have_sub: 1,
                        weight: dv.weight,
                    })
                    pushChildrenIData(dv.children, dv.filename)
                } else {
                    menuInfos.push({
                        menu_name: dv.menu_name,
                        is_top_menu: 1,
                        filename: dv.filename,
                        url: dv.url,
                        dir: "",
                        have_sub: 2,
                        weight: dv.weight,
                    })
                }
            }
        }
    }
    GetAllMenuInfo()
    console.log(menuInfos)
    return menuInfos
}