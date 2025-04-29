function GetAllFirstMenu() {
    const baseUrl = "%s"
    const lis = document.querySelector("ul.chunklist.chunklist_set > li")
    let menuInfos = []
    if (lis.length > 0) {
        lis.forEach(li => {
            const a = li.querySelector(":scope > a")
            if (a) {
                const menuName = a.textContent.trim()
                const url = a.href.trim()

                menuInfos.push({
                    menu_name: menuName
                })
            }


        })
    }
}