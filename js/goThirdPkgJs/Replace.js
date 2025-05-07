function agreeCookie() {
    const btn = document.querySelector(".Cookie-notice button.go-Button")
    if (btn) {
        btn.click()
    }
}

function expand() {
    const btn = document.querySelector("button.UnitReadme-expandLink")
    if (btn) {
        btn.click()
    }
}

function removeHeader() {
    const headers = document.querySelectorAll("header")
    if (headers.length > 0) {
        headers.forEach((h => {
            h.remove()
        }))
    }
}

function removeRenderedFor() {
    const ut = document.querySelector(".UnitBuildContext-titleContext")
    if (ut) {
        ut.remove()
    }
}

function removeDocumentationIndex() {
    const indexE = document.querySelector("section.Documentation-index")
    if (indexE) {
        indexE.remove()
    }
}

function removeDocumentationExamples() {
    const e = document.querySelector("section.Documentation-examples")
    if (e) {
        e.remove()
    }
}

function removeGoMainAside() {
    const aside = document.querySelector("aside.go-Main-aside")
    if (aside) {
        aside.remove()
    }
}

function removeGoNavigationDrawer() {
    const aside = document.querySelector("aside.go-NavigationDrawer")
    if (aside) {
        aside.remove()
    }
}

function removeGoMainNav() {
    const nav = document.querySelector("nav.go-Main-nav")
    if (nav) {
        nav.remove()
    }
}

function removeFooter() {
    const footer = document.querySelector("footer.go-Footer")
    if (footer) {
        footer.remove()
    }
}

function removeUnitFiles() {
    // UnitFiles js-unitFiles
    // UnitDirectories js-unitDirectories
    const divs = document.querySelectorAll("div.UnitFiles,div.UnitDirectories")
    if (divs.length > 0) {
        divs.forEach(div => {
            div.remove()
        })
    }
}

function removeUnitDirectories() {
    const div = document.querySelector("div.UnitDirectories")
    if (div) {
        div.remove()
    }
}

function removeImg() {
    const hs = document.querySelectorAll(".UnitDoc-title,.UnitReadme-title")
    if (hs.length > 0) {
        hs.forEach(h => {
            if (h.querySelector("img")) {
                h.querySelector("img").remove()
            }
        })
    }
}

function removeNoFollowImg() {
    const imgs = document.querySelectorAll(`a[rel="nofollow"] > img`)
    if (imgs.length > 0) {
        imgs.forEach(img => {
            img.remove()
        })
    }
}

// 打开 div.Documentation-deprecatedDetails
function openDeprecatedDetails() {
    const details = document.querySelectorAll(".Documentation-deprecatedDetails.js-deprecatedDetails")
    if (details.length > 0) {
        details.forEach(d => {
            d.setAttribute("open", "open")
            let sum = d.querySelector("summary")
            let subDiv = d.querySelector(":scope > div")

            while(sum && sum.firstElementChild) {
                d.parentElement.insertAdjacentElement("beforeend", sum.firstElementChild)
            }

            while(subDiv && subDiv.firstElementChild) {
                d.parentElement.insertAdjacentElement("beforeend", subDiv.firstElementChild)
            }
            d.remove()
        })
    }
}


// 往标题中修改锚
function addHeaderAnchorAndRemoveHeaderLink() {
    // 标题：Documentation
    const unitDocTitles = document.querySelectorAll(".UnitDoc-title")
    if (unitDocTitles.length > 0) {
        unitDocTitles.forEach(unitDocTitle => {
            const a = unitDocTitle.querySelector("a.UnitDoc-idLink")
            if (a) {
                const link = a.href
                //去除后面的 ¶
                const anchor = link.split("#")[1].trim()
                a.remove()
                unitDocTitle.insertAdjacentText("beforeend", `{#${anchor}}`)
            }
        })
    }

    // 标题：Overview、Constant、Variables、Functions、Types
    const docHeaders = document.querySelectorAll(".Documentation-overviewHeader,.Documentation-constantsHeader,.Documentation-variablesHeader,.Documentation-functionsHeader,.Documentation-typesHeader")
    if (docHeaders.length > 0) {
        docHeaders.forEach(docHeader => {
            const a = docHeader.querySelector("a")
            if (a) {
                const link = a.href
                //去除后面的 ¶
                const anchor = link.split("#")[1].trim()
                a.remove()
                docHeader.insertAdjacentText("beforeend", `{#${anchor}}`)
            }
        })
    }

    // 标题：函数名、类型名、方法名
    const ftts = document.querySelectorAll(".Documentation-functionHeader,.Documentation-typeHeader,.Documentation-typeMethodHeader,.Documentation-typeFuncHeader")
    if (ftts.length > 0) {
        ftts.forEach(ftt => {
            let level = ftt.tagName.replace("H", "")
            ftt.setAttribute("data-is-fmt", "yes")
            ftt.setAttribute("data-l", `${level}`)
            let versionStr = ""
            const outerSpan = ftt.querySelector("span.Documentation-sinceVersion")
            if (outerSpan) {
                // console.log("run hre")
                const innerSpan = outerSpan.querySelector("span.Documentation-sinceVersionVersion")
                if (innerSpan) {
                    versionStr = innerSpan.textContent.trim().replace("go", "")
                    versionStr = " <- " + versionStr
                }
                outerSpan.remove()
            }

            const a = ftt.querySelector("a.Documentation-idLink")
            // console.log("a=",a)
            if (a) {
                const link = a.href
                //去除后面的 ¶
                const anchor = link.split("#")[1].trim()
                a.remove()
                // console.log("ftt.firstElementChild=", ftt.firstElementChild)
                if (ftt.firstElementChild) {
                    ftt.firstElementChild.insertAdjacentText("beforeend", `${versionStr}{#${anchor}}`)
                }
            }
        })
    }

    // h2,h3,h4,h5中的a.Documentation-idLink
    const has = document.querySelectorAll("h2 > a.Documentation-idLink,h3 > a.Documentation-idLink,h4 > a.Documentation-idLink,h5 > a.Documentation-idLink,h6 > a.Documentation-idLink,.Documentation-notesHeader > a,.Documentation-noteHeader > a,.UnitReadme-title a.UnitReadme-idLink")
    if (has.length > 0) {
        has.forEach(ha => {
            const link = ha.href
            //去除后面的 ¶
            const anchor = link.split("#")[1].trim()
            const h = ha.parentElement
            ha.remove()
            h.insertAdjacentText("beforeend", `{#${anchor}}`)
        })
    }
}

// 修改在types下的函数和方法的标题层级
function modifyInTypeFuncHeaderLevel() {
    const fms = document.querySelectorAll(".Documentation-types .Documentation-typeMethodHeader,.Documentation-types .Documentation-typeFuncHeader")
    if (fms.length > 0) {
        fms.forEach(fm => {
            const newLevel = parseInt(fm.tagName.replace("H", "")) + 1
            const newH = document.createElement(`H${newLevel}`)
            newH.setAttribute("data-l", `${newLevel}`)
            newH.setAttribute("data-is-fmt", "yes")
            while (fm.firstChild) {
                newH.appendChild(fm.firstChild)
            }
            fm.insertAdjacentElement("afterend", newH)
            fm.remove()
        })
    }
}

// 替换Example中的代码块
function replaceExampleCodeBlock() {
    const des = document.querySelectorAll(".Documentation-exampleDetails")
    const noTopHLevel = 2
    const maxIterations = 100; // 最大迭代次数
    if (des.length > 0) {
        des.forEach(de => {
            const deb = de.querySelector(".Documentation-exampleDetailsBody")
            if (deb) {
                const ta = deb.querySelector(":scope > textarea")
                const pre = deb.querySelector(":scope > pre")
                let foundH = false
                let curMustSetHLevel = 3
                let anchor = ""
                // 找到最近h标签，并进行判断是否需要增加一级
                let prevE = de.previousElementSibling;
                // let deParent = de.parentElement
                // h标签中的标识符名称
                let idName = ""
                let iterationCount = 0; // 迭代计数器
                while (prevE && !foundH && iterationCount < maxIterations) {
                    if (['H1', 'H2', 'H3', 'H4', 'H5'].includes(prevE.tagName)) {
                        let dataIsExample = prevE.getAttribute("data-is-example")
                        if (dataIsExample) {
                            prevE = prevE.previousElementSibling;
                            if (!prevE) {
                                break
                            }
                            iterationCount++;
                            continue
                        }

                        let dataIsFmt = prevE.getAttribute("data-is-fmt")

                        if (dataIsFmt) {
                            foundH = true;
                            let dataL = prevE.getAttribute("data-l")
                            let idE = prevE.querySelector("a.Documentation-source")
                            if (idE) {
                                idName = idE.textContent.trim()
                            }
                            if (dataL) {
                                curMustSetHLevel = parseInt(dataL) + 1
                            }
                        }
                    } else {
                        if (prevE.previousElementSibling) {
                            prevE = prevE.previousElementSibling;
                        } else {
                            break
                        }
                    }
                    iterationCount++;
                }
                if (!foundH) {
                    curMustSetHLevel = noTopHLevel
                }

                // 获取Example原有锚
                const a = de.querySelector(".Documentation-exampleDetailsHeader a")
                if (a) {
                    anchor = a.href.split("#")[1].trim()
                }

                const div = document.createElement("div")
                const newH = document.createElement(`h${curMustSetHLevel}`)
                newH.textContent = `${idName} Example{#${anchor}}`
                newH.setAttribute("data-is-example", "yes")
                newH.setAttribute("data-l", `${curMustSetHLevel}`)
                div.appendChild(newH)
                // 获取代码块内容，并转换成html字符串
                div.insertAdjacentHTML("beforeend", convertCodeToHtml(ta.value, 'go'))
                // 插入一个换行
                div.insertAdjacentHTML("beforeend", "<br><p>&nbsp;&nbsp;&nbsp;&nbsp;Result：</p><br>")

                // 获取代码块的执行结果，并转换成html字符串
                div.insertAdjacentHTML("beforeend", convertCodeToHtml(pre.textContent, 'go'))

                de.insertAdjacentElement("afterend", div)
                de.remove()
            }
        })
    }
}

function convertCodeToHtml(code, lang) {
    // 将 \t 替换为四个空格
    let converted = code.replace(/\t/g, '&nbsp;&nbsp;&nbsp;&nbsp;');
    // 将 \n 替换为 <br> 标签
    converted = converted.replace(/\n\)$/, '<br><br>)');
    converted = converted.replace(/\n\}$/, '<br><br>}');
    converted = converted.replace(/\n/g, '<br>');
    // 使用 <pre> 和 <code> 标签包裹代码，并添加语言类名
    const html = `<div><pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-${lang}">${converted}</code></pre></div>`;
    return html;
}

// 替换已经存在的代码块
function replaceExistCodeBlock() {
    const pres = document.querySelectorAll(".Documentation-declaration > pre")
    if (pres.length > 0) {
        pres.forEach(pre => {
            pre.insertAdjacentHTML("afterend", convertCodeToHtml(pre.textContent, 'go'))
            pre.remove()
        })
    }
}

function replaceDocumentationDeprecatedTag() {
    const ddtags = document.querySelectorAll(".Documentation-deprecatedTag")
    if (ddtags.length > 0) {
        ddtags.forEach(ddtag => {
            ddtag.textContent = " <- DEPRECATED"
        })
    }
}



agreeCookie();
expand();
removeHeader();
removeRenderedFor();
removeFooter();
removeGoMainAside();
removeGoNavigationDrawer();
removeGoMainNav();
removeDocumentationIndex();
removeDocumentationExamples();
openDeprecatedDetails();
removeUnitFiles();
removeUnitDirectories();
removeImg();
removeNoFollowImg();
addHeaderAnchorAndRemoveHeaderLink();
modifyInTypeFuncHeaderLevel();
replaceDocumentationDeprecatedTag();
replaceExistCodeBlock();
replaceExampleCodeBlock();