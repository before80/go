// 移除不需要的元素 以及 替换复制不了的元素
// 移除div.related
function removeRelated() {
    document.querySelectorAll('div.related').forEach(div => {
        div.remove()
    });
}

// 移除页尾
function removeFooter() {
    document.querySelectorAll('div.footer').forEach(div => {
        div.remove()
    });
}

// 移除 H1
function removeH1() {
    document.querySelectorAll('h1').forEach(div => {
        div.remove()
    });
}

function removeSphinxsidebar() {
    document.querySelectorAll('div.sphinxsidebar').forEach(div => {
        div.remove()
    });
}

// 在em标签内容的左右两边加上反引号
function replaceEm() {
    document.querySelectorAll("em").forEach(em => {
        if (!['DT', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6'].includes(em.parentElement.tagName)) {
            if (em.classList.length === 0) {
                const span = document.createElement("span")
                span.innerHTML = "\u0060" + em.innerHTML + "\u0060";
                em.insertAdjacentElement("afterend", span)
                em.remove()
            }
        }
    })
}

function replaceEmStableabi() {
    document.querySelectorAll("em").forEach(em => {
        if (!['DT', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6'].includes(em.parentElement.tagName)) {
            if (em.classList.contains("stableabi")) {
                const pE = em.parentElement
                const p = document.createElement("p")
                p.appendChild(em)
                pE.insertAdjacentElement("afterbegin", p)
            }
        }
    })
}

// 在code标签内容的左右两边加上反引号
function replaceCode() {
    document.querySelectorAll("code").forEach(c => {
        if (!['A'].includes(c.parentElement.tagName)) {
            c.innerHTML = "\u0060" + c.innerHTML + "\u0060";
        }
    })
}

// 将div.versionadded 或 div.versionchanged 或div.admonition.note
// 或 div.admonition.seealso 或 div.admonition.caution
// 中的内容放在blockquote标签中
function replaceDivIntoBlockQuote() {
    document.querySelectorAll("div.versionadded,div.versionchanged,div.admonition").forEach(div => {
        const block = document.createElement("blockquote")
        while (div.firstChild) {
            block.appendChild(div.firstChild);
        }
        div.insertAdjacentElement("afterend", block)
        div.remove()
    })
}

// 将备注添加到其所指向的位置
function replaceRemarkToOriginPlace() {
    const footnoteListE = document.querySelector("aside.footnote-list")
    if (!footnoteListE) {
        return
    }
    const footnoteEs = footnoteListE.querySelectorAll("aside.footnote")
    footnoteEs.forEach(fn => {
        const content = fn.querySelector("p").textContent

        const aEs = fn.querySelectorAll(`span.label a[role="doc-backlink"]`)
        aEs.forEach(a => {
            const originId = a.href.split("#")[1].trim()
            console.log("originId=", originId)
            const xuHao = document.getElementById(originId).textContent
            const newContent = `（备注${xuHao}：${content}）`
            // 插入内容
            document.getElementById(originId).insertAdjacentText("afterend", newContent)
            document.getElementById(originId).remove()
        })
    })

    console.log(footnoteListE.previousElementSibling)
    if (footnoteListE.previousElementSibling.textContent.trim() === "备注") {
        footnoteListE.previousElementSibling.remove()
    }
    footnoteListE.remove()
}

// 去除代码块中 右上角的 >>>
function removeCopyButton() {
    document.querySelectorAll("span.copybutton").forEach(sp => {
        sp.remove()
    })
}

// 去除标题中固定链接（即一个悬浮图标）
function removeFixedHeaderLink() {
    document.querySelectorAll(".headerlink").forEach(a => {
        a.remove()
    })
}

// 增加添加标题的锚
function addHeaderAnchorAndRemoveHeaderLink() {
    const hs = document.querySelector("div.body").querySelectorAll("h2,h3,h4,h5,h6")

    hs.forEach(h => {
        const headerLinkE = h.querySelector("a.headerlink")
        if (headerLinkE) {
            const link = headerLinkE.href
            //去除后面的 ¶
            // const linkToTitle = h.textContent.replace(/¶+$/, '');
            const anchor = link.split("#")[1].trim()
            // h.setAttribute("data-href", link)
            headerLinkE.remove()
            // 判断下是否只有一个a标签了
            if (h.children.length === 1 && h.children[0].tagName === 'A') {
                const aTextContent = h.children[0].textContent
                h.children[0].remove()
                h.insertAdjacentText("beforeend", `${aTextContent}{#${anchor}}`)
            } else {
                h.insertAdjacentText("beforeend", `{#${anchor}}`)
            }
        }
    })
}

function setSectionDlLevel() {
    const dlEs = document.querySelectorAll("section > dl.py,section > dl.c")
    let sectionLevel = 2
    let foundH = false
    dlEs.forEach(dl => {
        if (!foundH) {
            let prevE = dl.previousElementSibling;
            while (prevE && !foundH) {
                if (['H1', 'H2', 'H3', 'H4', 'H5'].includes(prevE.tagName)) {
                    foundH = true;
                    sectionLevel = parseInt(prevE.tagName.replace("H", ""));
                } else {
                    prevE = prevE.previousElementSibling;
                }
            }
        }
    })

    dlEs.forEach(dl => {
        dl.setAttribute("data-h", `${sectionLevel + 1}`)
        dl.querySelectorAll(":scope > dd > dl").forEach(dl1 => {
            dl1.setAttribute("data-h", `${sectionLevel + 2}`)
            dl1.querySelectorAll(":scope > dd > dl").forEach(dl2 => {
                dl2.setAttribute("data-h", `${sectionLevel + 3}`)
                dl2.querySelectorAll(":scope > dd > dl").forEach(dl3 => {
                    dl3.setAttribute("data-h", `${sectionLevel + 4}`)
                })
            })
        })
    })
}

function replaceSectionDlDl(times) {
    console.log("times=", times)
    const dlEs = document.querySelectorAll("section > dl.py,section > dl.c")
    dlEs.forEach(dl => {
        const curMustSetToHLevel = dl.getAttribute('data-h')
        console.log("curMustSetToHLevel=",curMustSetToHLevel)
        const div = document.createElement('div');
        // 处理dl下的dt
        const dt = dl.querySelector(":scope > dt")
        if (dt) {
            const newH = document.createElement(`h${curMustSetToHLevel}`);
            let anchor = ""
            // 检查是否存在锚点
            if (dt.querySelector(".headerlink")) {
                const aE = dt.querySelector(".headerlink");
                anchor = aE.href.trim().split("#")[1];
                aE.remove(); // 移除原来的锚点链接
            }

            // while (dt.firstChild) {
            //     newH.appendChild(dt.firstChild)
            // }
            //
            // newH.insertBefore(document.createTextNode("`"), newH.firstChild)
            // newH.appendChild(document.createTextNode((anchor ? `{#${anchor}}` : "") + "`"))

            // 设置新标题的内容
            const title = dt.textContent.trim();
            newH.textContent = `\`${title}\`` + (anchor ? `{#${anchor}}` : "");
            div.appendChild(newH);
        }
        // 处理dl下的dd
        // 处理 dl 下的 dd
        const dd = dl.querySelector(":scope > dd");
        if (dd) {
            // 将 dd 的子节点移动到新的 div 容器中
            while (dd.firstChild) {
                div.appendChild(dd.firstChild);
            }
        }

        let child = div.firstChild;
        while (child) {
            const nextChild = child.nextSibling;
            if (child.nodeType === Node.ELEMENT_NODE) {
                dl.insertAdjacentElement('beforebegin', child);
            }
            child = nextChild;
        }
        dl.remove()
    })
}

function repeatReplaceSectionDl() {
    setSectionDlLevel()
    replaceSectionDlDl(1)
    replaceSectionDlDl(2)
    replaceSectionDlDl(3)
    // replaceSectionDlDl()
    // replaceSectionDlDl()
}



// 修改并替换dl.py.method 或 dl.py.attribute或dl.c.function 或dl.c.var或dl.c.macro中的内容
function replaceDlPyC() {
    const dlEs = document.querySelectorAll("dl.py,dl.c")
    dlEs.forEach(dl => {
        //找到最近的标题是h3还是h4还是h5
        let prevE = dl.previousElementSibling;
        let foundH = false
        let curMustSetToHLevel = 4
        while (prevE && !foundH) {
            if (['H1', 'H2', 'H3', 'H4', 'H5'].includes(prevE.tagName)) {
                foundH = true;
                curMustSetToHLevel = parseInt(prevE.tagName.replace("H", "")) + 1;
            } else {
                prevE = prevE.previousElementSibling;
            }
        }

        const fragment = document.createDocumentFragment();
        // 处理dl下的dt
        const dts = dl.querySelectorAll("dt")
        dts.forEach(dt => {
            const divForDt = document.createElement("div")
            const newH = document.createElement(`h${curMustSetToHLevel}`);
            let anchor = ""
            // 检查是否存在锚点
            if (dt.querySelector(".headerlink")) {
                const aE = dt.querySelector(".headerlink");
                anchor = aE.href.trim().split("#")[1];
                aE.remove(); // 移除原来的锚点链接
            }

            // while (dt.firstChild) {
            //     newH.appendChild(dt.firstChild)
            // }

            // newH.insertBefore(document.createTextNode("`"), newH.firstChild)
            // newH.appendChild(document.createTextNode((anchor ? `{#${anchor}}` : "") + "`"))
            // 设置新标题的内容
            const title = dt.textContent.trim();
            newH.textContent = `\`${title}\`` + (anchor ? `{#${anchor}}` : "");
            divForDt.appendChild(newH)
            fragment.appendChild(divForDt);
        })

        // 处理dl下的dd
        // 处理 dl 下的 dd
        const dd = dl.querySelector("dd");
        if (dd) {
            const divForDd = document.createElement("div");
            // 将 dd 的子节点移动到新的 div 容器中
            while (dd.firstChild) {
                divForDd.appendChild(dd.firstChild);
            }
            fragment.appendChild(divForDd);
        }

        while (fragment.lastChild) {
            dl.insertAdjacentElement('afterend', fragment.lastChild);
        }
        dl.remove()
    })
}

function replaceDlGlossary() {
    const dls = document.querySelectorAll("dl.glossary")
    if (dls.length > 0) {
        dls.forEach(dl => {
            const fragment = document.createDocumentFragment();
            let curMustSetToHLevel = 2
            let prevE = dl.previousElementSibling;
            let foundH = false
            while (prevE && !foundH) {
                if (['H1', 'H2', 'H3', 'H4', 'H5'].includes(prevE.tagName)) {
                    foundH = true;
                    curMustSetToHLevel = parseInt(prevE.tagName.replace("H", "")) + 1;
                } else {
                    prevE = prevE.previousElementSibling;
                }
            }

            const dts = dl.querySelectorAll(":scope > dt")
            if (dts.length > 0) {
                dts.forEach(dt => {
                    const newH = document.createElement(`h${curMustSetToHLevel}`);
                    let anchor = ""
                    // 检查是否存在锚点
                    if (dt.querySelector(".headerlink")) {
                        const aE = dt.querySelector(".headerlink");
                        anchor = aE.href.trim().split("#")[1];
                        aE.remove(); // 移除原来的锚点链接
                    }
                    // 设置新标题的内容
                    const title = dt.textContent.trim();
                    newH.textContent = `\`${title}\`` + (anchor ? `{#${anchor}}` : "");
                    fragment.appendChild(newH)

                    const dd = dt.nextElementSibling
                    if (dd) {
                        const divForDd = document.createElement("div");
                        // 将 dd 的子节点移动到新的 div 容器中
                        while (dd.firstChild) {
                            divForDd.appendChild(dd.firstChild);
                        }
                        fragment.appendChild(divForDd);
                    }
                })
            }
            while (fragment.lastChild) {
                dl.insertAdjacentElement('afterend', fragment.lastChild);
            }
            dl.remove()
        })
    }
}

function replaceIntoBlockquote() {
    document.querySelectorAll("aside.topic").forEach(aside => {
        const blockquote = document.createElement('blockquote')
        while (aside.firstChild) {
            blockquote.appendChild(aside.firstChild);
        }
        aside.insertAdjacentElement('afterend', blockquote);
        aside.remove()
    })

    document.querySelectorAll("div.deprecated").forEach(div => {
        const blockquote = document.createElement('blockquote')
        while (div.firstChild) {
            blockquote.appendChild(div.firstChild);
        }
        div.insertAdjacentElement('afterend', blockquote);
        div.remove()
    })

    document.querySelectorAll("div.deprecated-removed").forEach(div => {
        const blockquote = document.createElement('blockquote')
        while (div.firstChild) {
            blockquote.appendChild(div.firstChild);
        }
        div.insertAdjacentElement('afterend', blockquote);
        div.remove()
    })

}

// 替换 dl.field-list.simple 中的内容
function replaceDlFieldList() {
    document.querySelectorAll('dl.field-list.simple').forEach(dl => {
        const blockquote = document.createElement('blockquote')
        const dt = dl.querySelector('dt')
        if (dt) {
            const divForDt = document.createElement('div')
            while (dt.firstChild) {
                divForDt.appendChild(dt.firstChild);
            }
            blockquote.appendChild(divForDt)
        }

        const dds = dl.querySelectorAll('dd')
        if (dds.length > 0) {
            dds.forEach(dd => {
                const divForDd = document.createElement('div')
                while (dd.firstChild) {
                    divForDd.appendChild(dd.firstChild);
                }
                blockquote.appendChild(divForDd)
            })
        }

        dl.insertAdjacentElement('afterend', blockquote);
        dl.remove()
    })
}

// 在p标签前面添加一个span，其内容为&zeroWidthSpace;用于后期在markdown中替换成Tab符号
function replaceP() {
    // 在p标签前面插入&zeroWidthSpace;
    document.querySelectorAll('p').forEach(function (p) {
        if (!['LI', "BLOCKQUOTE", "TH", "TD"].includes(p.parentElement.tagName)) {
            let newSpan = document.createElement('span');
            newSpan.textContent = '&zeroWidthSpace;';
            if (p.firstChild) {
                p.insertBefore(newSpan, p.firstChild);
            } else {
                // 如果 p 元素没有子节点，直接将新 span 元素添加到 p 元素中
                p.appendChild(newSpan);
            }
        }
    });
}

function removeTocTree() {
    const tocTree = document.querySelector("div.toctree-wrapper.compound")
    if (tocTree) {
        tocTree.remove()
    }

    const navContents = document.querySelector("nav.contents")
    if (navContents) {
        navContents.remove()
    }
}

removeRelated()
removeFooter();
removeH1();
removeSphinxsidebar();
replaceEm();
replaceEmStableabi();
// replaceCode();
replaceDivIntoBlockQuote();
replaceRemarkToOriginPlace();
removeCopyButton();
replaceDlGlossary();
addHeaderAnchorAndRemoveHeaderLink();
replaceIntoBlockquote();
repeatReplaceSectionDl();
replaceDlPyC();
replaceDlFieldList();
replaceP();
removeTocTree()