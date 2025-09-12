// The generator view script
(function () {
  function main() {
    prepareTemplateBuilder();
    prepareCitationGenerator();
  }

  // --- Template builder ---
  // Builds template from user defined fields

  // WeakMap holding custom properties and methods for field elements to avoid storing them in the HTMLElements themselfs
  const fieldCustomData = new WeakMap();

  function prepareTemplateBuilder() {
    const templateBuilder = document.getElementById("builder");
    if (templateBuilder === null) {
      throw new Error("Element with id 'builder' must exist");
    }
    templateBuilder.addEventListener("dragover", (e) => {
      e.preventDefault();
    });
    templateBuilder.addEventListener("drop", (e) =>
      dropHanlder(e, templateBuilder)
    );

    const fieldTypeSelect = document.getElementById("field-type");
    if (fieldTypeSelect === null) {
      throw new Error("Element with id 'field-type' must exist");
    }

    const addFieldBtn = document.getElementById("add-field");
    if (addFieldBtn === null) {
      throw new Error("Element with id 'add-field' must exist");
    }
    let fieldNumber = 1;
    addFieldBtn.addEventListener("click", () => {
      addField(templateBuilder, fieldTypeSelect.value, fieldNumber);
      fieldNumber++;
    });

    const removeAllFieldsBtn = document.getElementById("remove-all-fields");
    if (removeAllFieldsBtn === null) {
      throw new Error("Element with id 'remove-all-fields' must exist");
    }
    removeAllFieldsBtn.addEventListener(
      "click",
      () => (templateBuilder.innerHTML = "")
    );

    const templateInputElement = document.getElementById("template");
    if (templateInputElement === null) {
      throw new Error("Element with id 'template' must exist");
    }
    const buildTemplateBtn = document.getElementById("build-template");
    if (buildTemplateBtn === null) {
      throw new Error("Element with id 'build-template' must exist");
    }
    buildTemplateBtn.addEventListener("click", () => {
      buildTemplate(templateBuilder, templateInputElement);
      templateInputElement.dispatchEvent(new Event("input", { bubbles: true }));
    });
  }

  /**
   * @param {DragEvent} e
   * @param {HTMLElement} templateBuilder
   */
  function dropHanlder(e, templateBuilder) {
    e.preventDefault();
    const fieldCount = templateBuilder.childElementCount;
    if (fieldCount <= 1) {
      // Only one field which must be the origin. Return early.
      return;
    }
    const originID = e.dataTransfer.getData("text/plain");
    const originField = document.getElementById(originID);
    const currentY = e.clientY;
    const fields = templateBuilder.children;
    for (let i = 0; i < fieldCount; i++) {
      const field = fields[i];
      const rect = field.getBoundingClientRect();
      if (currentY <= rect.y + rect.height / 2) {
        field.insertAdjacentElement("beforebegin", originField);
        return; // Done
      }
    }
    // If nothing matched, put it at the end
    templateBuilder.insertAdjacentElement("beforeend", originField);
  }

  /**
   * @param {HTMLElement} target The element under which the new field will be added
   * @param {string} type
   * @param {number} fieldNumber
   */
  function addField(target, type, fieldNumber) {
    const id = "field-" + fieldNumber;
    let fieldInitFunc = null;
    switch (type) {
      case "text":
        fieldInitFunc = initTextField;
        break;
      case "jméno":
        fieldInitFunc = initNameField;
        break;
      case "příjmení":
        fieldInitFunc = initLastnameField;
        break;
      case "název":
        fieldInitFunc = initWebNameField;
        break;
      case "součást":
        fieldInitFunc = initPartOfField;
        break;
      case "místo-vydání":
        fieldInitFunc = initPlaceOfPublicationField;
        break;
      case "datum-vydání":
        fieldInitFunc = initDateOfPublicationField;
        break;
      case "url":
        fieldInitFunc = initUrlField;
        break;
      case "archivní-url":
        fieldInitFunc = initArchivalUrlField;
        break;
      case "datum-archivace":
        fieldInitFunc = initDateOfHarvestField;
        break;
      case "datum-citace":
        fieldInitFunc = initDatefCitationField;
        break;
      default:
        throw new Error(`unknown field type: ${type}`);
    }
    target.append(createNewField(type, id, fieldInitFunc));
  }

  /**
   * @param {string} type
   * @param {string} id
   * @param {function(HTMLFormElement): void} initField
   * @returns {HTMLFormElement}
   */
  function createNewField(type, id, initField) {
    const field = document.createElement("form");
    field.id = id;
    field.classList.add("flex-row");
    field.classList.add("field");

    field.dataset.type = type;

    // Add the field to the WeakMap so the "getTemplateValue" method can be safely stored
    fieldCustomData.set(field, {});

    initField(field);

    const removeBtn = document.createElement("button");
    removeBtn.append("Odebrat");
    removeBtn.addEventListener("click", () => field.remove());
    field.append(removeBtn);

    field.draggable = true;
    field.addEventListener("dragstart", (e) => {
      e.dataTransfer.setData("text/plain", field.id);
      e.dataTransfer.dropEffect = "move";
    });

    return field;
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initTextField(field) {
    field.innerHTML = `
      <span class="f-start">Textové&nbsp;pole:</span>
      <div class="flex-row max-flex f-middle">
        <label class="flex-row max-flex"><input class="max-flex" type="text" name="f-value"></label>
        ${fieldFormatFormControls}
      </div>
    `;
    stopElementFromBeingDragged(field.elements.namedItem("f-value"), field);
    const data = fieldCustomData.get(field);
    data.getTemplateValue = function () {
      const text = field.elements.namedItem("f-value").value;
      return wrapTextInFormat(text, field);
    };
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initNameField(field) {
    initGenericFormatAndSeparatorField(field, "Jméno&nbsp;autora", "jméno");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initLastnameField(field) {
    initGenericFormatAndSeparatorField(
      field,
      "Příjmení&nbsp;autora",
      "příjmení"
    );
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initWebNameField(field) {
    initGenericFormatAndSeparatorField(field, "Název&nbsp;webu", "název");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initPartOfField(field) {
    initGenericFormatAndSeparatorField(field, "Součást", "součást");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initPlaceOfPublicationField(field) {
    initGenericFormatAndSeparatorField(
      field,
      "Místo&nbsp;vydání",
      "místo-vydání"
    );
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initDateOfPublicationField(field) {
    initGenericFormatAndSeparatorField(
      field,
      "Datum&nbsp;vydání",
      "datum-vydání"
    );
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initUrlField(field) {
    initGenericFormatAndSeparatorField(field, "URL", "url");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initArchivalUrlField(field) {
    initGenericFormatAndSeparatorField(
      field,
      "Archivní&nbsp;URL",
      "archivní-url"
    );
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initDateOfHarvestField(field) {
    initGenericFormatAndSeparatorField(
      field,
      "Datum archivace",
      "datum-archivace"
    );
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initDatefCitationField(field) {
    initGenericFormatAndSeparatorField(field, "Datum citace", "datum-citace");
  }

  /**
   * @param {HTMLFormElement} field
   * @param {string} readableName
   * @param {string} exprName
   */
  function initGenericFormatAndSeparatorField(field, readableName, exprName) {
    field.innerHTML = `
      <span class="f-start">${readableName}:</span>
      <div class="flex-column max-flex f-middle">
        <div class="flex-row max-flex">
          ${fieldFormatFormControls}
          ${fieldCaseFormControls}
        </div>
        <div class="flex-row">
          ${fieldSeparatorFormControls}
        </div>
      </div>
    `;
    stopElementFromBeingDragged(field.elements.namedItem("f-oddělovač"), field);
    const data = fieldCustomData.get(field);
    data.getTemplateValue = function () {
      let expr = exprName;
      expr = setCase(expr, field);
      expr = wrapExprInFormat(expr, field);
      expr = `{{${expr}}}`;
      expr = addSeparator(expr, field);
      return expr;
    };
  }

  /**
   * This is complete hack but it kinda works.
   * Use this on text input elements to allow selecting text etc.
   * It sometimes leaves the field in undragabble state but it's rare and clicking it will fix it.
   * @param {HTMLElement} element
   * @param {HTMLElement} field
   */
  function stopElementFromBeingDragged(element, field) {
    element.addEventListener("mousedown", (e) => (field.draggable = false));
    element.addEventListener("mouseup", (e) => (field.draggable = true));
    field.addEventListener("mouseup", (e) => (field.draggable = true));
  }

  const fieldSeparatorFormControls = `
    <label class="flex-row">Interpunkce:<input type="text" name="f-oddělovač"></label>
  `;
  /**
   * Will append separator at the end of expr. Use as last function in chain.
   * @param {string} expr
   * @param {HTMLFormElement} field
   * @return {string}
   */
  function addSeparator(expr, field) {
    const separator = field.elements.namedItem("f-oddělovač")?.value;
    if (!separator) {
      return expr;
    }
    return expr + separator;
  }

  const fieldFormatFormControls = `
  <span class="flex-row format-controls">
    <label class="flex-row"><input type="checkbox" name="f-tučně">Tučně</label>
    <label class="flex-row"><input type="checkbox" name="f-kurzíva">Kurzívou</label>
  </span>
  `;
  /**
   * Add formating to handlebars expression.
   * @param {string} expr
   * @param {HTMLFormElement} field
   * @return {string}
   */
  function wrapExprInFormat(expr, field) {
    if (field.elements.namedItem("f-tučně")?.checked) {
      if (expr.split(" ").length > 1) {
        expr = `tučně (${expr})`;
      } else {
        expr = `tučně ${expr}`;
      }
    }
    if (field.elements.namedItem("f-kurzíva")?.checked) {
      if (expr.split(" ").length > 1) {
        expr = `kurzíva (${expr})`;
      } else {
        expr = `kurzíva ${expr}`;
      }
    }
    return expr;
  }
  /**
   * Add formating to plain text.
   * @param {string} text
   * @param {HTMLFormElement} field
   * @return {string}
   */
  function wrapTextInFormat(text, field) {
    if (field.elements.namedItem("f-tučně")?.checked) {
      text = "{{tučně}}" + text + "{{-tučně}}";
    }
    if (field.elements.namedItem("f-kurzíva")?.checked) {
      text = "{{kurzíva}}" + text + "{{-kurzíva}}";
    }
    return text;
  }

  const fieldCaseFormControls = `
  <span class="flex-row case-controls">
    <label class="flex-row"><input type="radio" name="f-case" value="no-change" checked>Neměnit</label>
    <label class="flex-row"><input type="radio" name="f-case" value="small">Malé</label>
    <label class="flex-row"><input type="radio" name="f-case" value="capital">Velké</label>
    <label class="flex-row"><input type="radio" name="f-case" value="first-capital">První&nbsp;velké</label>
  </span>
  `;
  /**
   * @param {string} expr
   * @param {HTMLFormElement} field
   * @return {string}
   */
  function setCase(expr, field) {
    let helperName = "";
    switch (field.elements.namedItem("f-case").value) {
      case "small":
        helperName = "malé";
        break;
      case "capital":
        helperName = "verzálky";
        break;
      case "first-capital":
        helperName = "první-velké";
        break;
      default:
        return expr;
    }
    if (expr.split(" ").length > 1) {
      expr = `${helperName} (${expr})`;
    } else {
      expr = `${helperName} ${expr}`;
    }
    return expr;
  }

  /**
   *
   * @param {HTMLElement} templateBuilder
   * @param {HTMLElement} target
   */
  function buildTemplate(templateBuilder, target) {
    let template = "";
    for (const field of templateBuilder.children) {
      const data = fieldCustomData.get(field);
      if (!data?.getTemplateValue) {
        continue;
      }
      template += data.getTemplateValue() + " ";
    }
    target.value = template;
  }

  // --- Citation generator ---
  // Generates citations from template

  function prepareCitationGenerator() {
    // Get elements for input form and output paragraph.
    const generatorForm = document.getElementById("generator");
    if (!(generatorForm instanceof HTMLFormElement)) {
      throw new Error(
        "Element with id 'generator' must exist and be HTMLFormElement"
      );
    }

    const templateElement = document.getElementById("template");
    if (templateElement === null) {
      throw new Error("Element with id 'template' must exist");
    }

    const citationOutput = document.getElementById("citation");
    if (citationOutput === null) {
      throw new Error("Element with id 'citation' must exist");
    }

    const inputDataElement = document.getElementById("input-data");
    // If input data are present parse them and fill the Form.
    if (inputDataElement !== null) {
      try {
        if (inputDataElement.type !== "application/json") {
          throw new Error(
            "Element #input-data must be script with type 'application/json'"
          );
        }
        const inputData = JSON.parse(inputDataElement.text);
        if (!Array.isArray(inputData)) {
          throw new Error("inputData must be array of objects");
        }
        if (inputData.length !== 0) {
          fillForm(generatorForm, inputData[0]);
          enableFormControls(
            generatorForm,
            templateElement,
            citationOutput,
            inputData
          );
        }
      } catch (err) {
        console.error(err);
      }
    }

    // Render the template any time when user inputs data.
    generatorForm.addEventListener("input", () =>
      generateCitation(generatorForm, templateElement, citationOutput)
    );
    // Also when template is changed
    templateElement.addEventListener("input", () =>
      generateCitation(generatorForm, templateElement, citationOutput)
    );

    // Render the template first time on page load.
    generateCitation(generatorForm, templateElement, citationOutput);
  }

  /**
   *
   * @param {HTMLFormElement} generatorForm
   * @param {HTMLElement} templateElement
   * @param {HTMLElement} citationOutput
   * @param {Array} citationData
   */
  function enableFormControls(
    generatorForm,
    templateElement,
    citationOutput,
    citationData
  ) {
    const formControls = document.getElementById("form-controls");
    if (formControls === null) {
      throw new Error("Element with id 'form-controls' must exist");
    }

    formControls.hidden = false;
    formControls.classList.remove("hidden");

    let currentDataIndex = 0;
    let dataCount = citationData.length;

    const currentIndexElem = document.getElementById("cit-data-num");
    const countElem = document.getElementById("cit-data-count");
    if (currentIndexElem === null) {
      throw new Error("Element with id 'cit-data-num' must exist");
    }
    if (countElem === null) {
      throw new Error("Element with id 'cit-data-count' must exist");
    }

    // Show index bigger by one as that is what people generaly expect
    currentIndexElem.textContent = (currentDataIndex + 1).toString();
    countElem.textContent = dataCount.toString();

    const prevBtn = document.getElementById("prev");
    const nextBtn = document.getElementById("next");
    if (prevBtn === null) {
      throw new Error("Element with id 'prev' must exist");
    }
    if (nextBtn === null) {
      throw new Error("Element with id 'next' must exist");
    }

    prevBtn.addEventListener("click", () => {
      if (currentDataIndex <= 0) {
        return;
      }
      saveCurrentCitationData(generatorForm, citationData, currentDataIndex);
      currentDataIndex--;
      currentIndexElem.textContent = (currentDataIndex + 1).toString();
      fillForm(generatorForm, citationData[currentDataIndex]);
      generateCitation(generatorForm, templateElement, citationOutput);
    });

    nextBtn.addEventListener("click", () => {
      if (currentDataIndex >= dataCount - 1) {
        return;
      }
      saveCurrentCitationData(generatorForm, citationData, currentDataIndex);
      currentDataIndex++;
      currentIndexElem.textContent = (currentDataIndex + 1).toString();
      fillForm(generatorForm, citationData[currentDataIndex]);
      generateCitation(generatorForm, templateElement, citationOutput);
    });
  }

  /**
   *
   * @param {HTMLFormElement} generatorForm
   * @param {Array} citationData
   * @param {number} currentIndex
   */
  function saveCurrentCitationData(generatorForm, citationData, currentIndex) {
    const currentData = citationData[currentIndex];
    for (const element of generatorForm.elements) {
      if ("citationfield" in element.dataset) {
        currentData[element.dataset.citationfield] = element.value;
      }
    }
  }

  /**
   *
   * @param {HTMLFormElement} generatorForm
   * @param {any} inputData // Object containing citation field values
   */
  function fillForm(generatorForm, inputData) {
    for (const field in inputData) {
      const element = generatorForm.elements.namedItem(field);
      if (element === null) {
        console.warn(`Form element with id ${field} does not exist!`);
        continue;
      }
      element.value = inputData[field];
    }
  }

  /**
   * Replace citationOutput element with template and data from generatorForm.
   * @param {HTMLFormElement} generatorForm
   * @param {HTMLElement} templateElement
   * @param {HTMLElement} citationOutput
   */
  function generateCitation(generatorForm, templateElement, citationOutput) {
    const template = Handlebars.compile(templateElement.value);
    const data = {};
    for (const element of generatorForm.elements) {
      if ("citationfield" in element.dataset) {
        data[element.dataset.citationfield] = element.value;
      }
    }
    const citation = template(data);
    citationOutput.innerHTML = citation;
    console.log(citation);
  }

  // Add helpers to Handlebars for producing formatted output.
  Handlebars.registerHelper("tučně", (text) => formatBold(false, text));
  Handlebars.registerHelper("b", (text) => formatBold(false, text));
  Handlebars.registerHelper("-tučně", (text) => formatBold(true, text));
  Handlebars.registerHelper("-b", (text) => formatBold(true, text));
  Handlebars.registerHelper("kurzíva", (text) => formatItalic(false, text));
  Handlebars.registerHelper("i", (text) => formatItalic(false, text));
  Handlebars.registerHelper("-kurzíva", (text) => formatItalic(true, text));
  Handlebars.registerHelper("-i", (text) => formatItalic(true, text));
  Handlebars.registerHelper("verzálky", upperCase);
  Handlebars.registerHelper("první-velké", capitalize);
  Handlebars.registerHelper("malé", lowerCase);

  /**
   * @param {boolean} end
   * @param {any} text
   */
  function formatBold(end, text) {
    console.log(text);
    const startingTag = "<b>";
    const endingTag = "</b>";
    if (end) {
      return new Handlebars.SafeString(endingTag);
    }
    if (text instanceof Handlebars.SafeString) {
      return new Handlebars.SafeString(startingTag + text.string + endingTag);
    }
    if (!text || typeof text !== "string") {
      return new Handlebars.SafeString(startingTag);
    }
    text = Handlebars.escapeExpression(text);
    return new Handlebars.SafeString(startingTag + text + endingTag);
  }

  /**
   * @param {boolean} end
   * @param {any} text
   */
  function formatItalic(end, text) {
    console.log(text);
    const startingTag = "<i>";
    const endingTag = "</i>";
    if (end) {
      return new Handlebars.SafeString(endingTag);
    }
    if (text instanceof Handlebars.SafeString) {
      return new Handlebars.SafeString(startingTag + text.string + endingTag);
    }
    if (!text || typeof text !== "string") {
      return new Handlebars.SafeString(startingTag);
    }
    text = Handlebars.escapeExpression(text);
    return new Handlebars.SafeString(startingTag + text + endingTag);
  }

  /**
   * @param {string} text
   * @returns {string}
   */
  function upperCase(text) {
    return text.toUpperCase();
  }

  /**
   * @param {string} text
   * @returns {string}
   */
  function capitalize(text) {
    return text
      .split(" ")
      .map((word) => word[0].toUpperCase() + word.slice(1))
      .join(" ");
  }

  /**
   * @param {string} text
   * @returns {string}
   */
  function lowerCase(text) {
    return text.toLowerCase();
  }

  main();
})();
