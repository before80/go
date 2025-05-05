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
}

function replaceSomething() {
    const tdSpanStrongs = document.querySelectorAll("td > span.bold > strong")
    if (tdSpanStrongs.length > 0) {
        tdSpanStrongs.forEach(st => {
            const span = st.parentElement
            const newSt = document.createElement("strong")
            newSt.innerHTML = st.innerHTML
            span.insertAdjacentElement("afterend", newSt)
            span.remove()
        })
    }
}

removeSomething();
replaceSomething();