() => {
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
                        menu_name: a.textContent.trim(),
                        filename: filename,
                        url: url,
                        index: i + 1,
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

    function pushChildrenIData(menuInfos, chd, parentDirPath) {
        const iData = chd
        for (let i = 0; i < iData.length; i++) {
            const dv = iData[i]
            let dirPath = dv.children && dv.children.length > 0 ? `${parentDirPath}/{dv.filename}` : parentDirPath
            menuInfos.push({
                menu_name: dv.menu_name,
                filename: dv.filename,
                file_path: dv.children && dv.children.length > 0 ? `${dirPath}/${dv.filename}/_index.md` : `${dirPath}/${dv.filename}.md`,
                index: dv.index,
                isTop: 2,
                children: dv.children,
            })
            if (dv.children && dv.children.length > 0) {
                pushChildrenIData(dv.children, dirPath)
            }
        }
    }

    function GetAllMenuInfo() {
        let tempMenuInfos = []
        tempMenuInfos = getNestedLiData(document.querySelector("#doc-201 > ul"))

        let menuInfos = []

        if (tempMenuInfos.length > 0) {
            for (let i = 0; i < tempMenuInfos.length; i++) {
                const dv = tempMenuInfos[i]
                let dirPath =  dv.filename
                menuInfos.push({
                    menu_name: dv.menu_name,
                    filename: dv.filename,
                    file_path: `${dirPath}/_index.md`,
                    index: dv.index,
                    is_top: 1,
                })
                if (dv.children && dv.children.length > 0) {
                    pushChildrenIData(menuInfos, dv.children, dirPath)
                }
            }
        }
        return menuInfos
    }

    return GetAllMenuInfo()
}