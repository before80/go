

function removeHeader() {
    const headers = document.querySelectorAll("header")
    if (headers.length > 0) {
        headers.forEach((h => {
            h.remove()
        }))
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
    const div = document.querySelector("div.UnitFiles")
    if (div) {
        div.remove()
    }
}

function removeUnitDirectories() {
    const div = document.querySelector("div.UnitDirectories")
    if (div) {
        div.remove()
    }
}

function replaceUnitDocTitle() {
    const h = document.querySelector(".UnitDoc-title")
    if (h) {
        if (h.querySelector("img")) {
            h.querySelector("img").remove()
        }
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
}

// 修改在types下的函数和方法的标题层级
function modifyInTypeFuncHeaderLevel() {
    const fms = document.querySelectorAll(".Documentation-types .Documentation-typeMethodHeader,.Documentation-types .Documentation-typeFuncHeader")
    if (fms.length > 0) {
        fms.forEach(fm => {
            const newLevel = parseInt(fm.tagName.replace("H", "")) + 1
            const newH = document.createElement(`H${newLevel}`)
            newH.setAttribute("data-h", `${newLevel}`)
            newH.setAttribute("data-is-new-header", "no")
            while(fm.firstChild) {
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
    if (des.length > 0) {
        des.forEach(de => {
            const ta = de.querySelector(".Documentation-exampleDetailsBody > textarea")
            if (ta) {
                // 找到最近h标签，并进行判断是否需要增加一级

                // 获取h标签中的标识符名称

            }
        })
    }
}

function convertCodeToHtml(code,lang) {
    // 将 \t 替换为四个空格
    let converted = code.replace(/\t/g, '&nbsp;&nbsp;&nbsp;&nbsp;');
    // 将 \n 替换为 <br> 标签
    converted = converted.replace(/\n\)$/, '<br><br>)');
    converted = converted.replace(/\n\}$/, '<br><br>}');
    converted = converted.replace(/\n/g, '<br>');
    // 使用 <pre> 和 <code> 标签包裹代码，并添加语言类名
    const html = `<pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-${lang}">${converted}</code></pre>`;
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




removeHeader();
removeFooter();
removeGoMainAside();
removeGoMainNav();
removeDocumentationIndex();
removeDocumentationExamples();
removeUnitFiles();
removeUnitDirectories();
replaceUnitDocTitle();
replaceExampleCodeBlock();
replaceExistCodeBlock();
addHeaderAnchorAndRemoveHeaderLink();
modifyInTypeFuncHeaderLevel();
