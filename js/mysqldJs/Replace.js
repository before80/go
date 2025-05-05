function removeSomething() {
    const header = document.querySelector("header")
    if (header) {
        header.remove()
    }

    const footer = document.querySelector("footer")
    if (footer) {
        footer.remove()
    }

    const sidebar = document.querySelector("#docs-sidebar-toc")
    if (sidebar) {
        sidebar.remove()
    }

    const skipMain = document.querySelector("a.skip-main")
    if (skipMain) {
        skipMain.remove()
    }

    const version = document.querySelector("#docs-version-nav")
    if (version) {
        version.remove()
    }

    const show = document.querySelector("#docs-show-nav")
    if (show) {
        show.remove()
    }

    const breadcrumbs = document.querySelector("#docs-breadcrumbs")
    if (breadcrumbs) {
        breadcrumbs.remove()
    }

    const container = document.querySelector("#docs-in-page-nav-container")
    if (container) {
        container.remove()
    }

    const doc = document.querySelector("div.toc")
    if (doc) {
        doc.remove()
    }

    const indexterms = document.querySelectorAll("a.indexterm")
    if (indexterms.length > 0) {
        indexterms.forEach(it => {
            it.remove()
        })
    }

}

function replaceSomething() {
    const spanStrongs = document.querySelectorAll("span.bold > strong")
    if (spanStrongs.length > 0) {
        spanStrongs.forEach(st => {
            const span = st.parentElement
            const newSt = document.createElement("strong")
            newSt.innerHTML = st.innerHTML
            span.insertAdjacentElement("afterend", newSt)
            span.remove()
        })
    }

    const literalLinks = document.querySelectorAll("code.literal > a.link")

    if (literalLinks.length > 0) {
        literalLinks.forEach(ll => {
            const code = ll.parentElement
            code.insertAdjacentElement("afterend", ll)
            code.remove()
        })
    }

    const emCodes = document.querySelectorAll("em.replaceable > code")

    if (emCodes.length > 0) {
        emCodes.forEach(ec => {
            const em = ec.parentElement
            em.insertAdjacentElement("afterend", ec)
            em.remove()
        })
    }

}


function addBlockquote() {
    const divs = document.querySelectorAll("div.note,div.important,div.tip,div.warning")
    if (divs.length > 0) {
        divs.forEach(div => {
            const blockquote = document.createElement('blockquote')
            const newDiv = document.createElement('div')
            newDiv.insertAdjacentHTML("afterbegin", div.innerHTML)
            blockquote.appendChild(newDiv)
            div.insertAdjacentElement("afterend", blockquote)
            div.remove()
        })
    }
}

removeSomething();
replaceSomething();
addBlockquote()