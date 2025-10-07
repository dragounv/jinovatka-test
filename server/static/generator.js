// The generator view script
// This script needs Handlebars and Luxon to be loaded
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
      // Remove placeholder if present
      const placeholder = document.getElementById("builder-placeholder");
      if (placeholder) {
        placeholder.remove();
      }
    });

    const removeAllFieldsBtn = document.getElementById("remove-all-fields");
    if (removeAllFieldsBtn === null) {
      throw new Error("Element with id 'remove-all-fields' must exist");
    }
    removeAllFieldsBtn.addEventListener(
      "click",
      () =>
        (templateBuilder.innerHTML = `<i id="builder-placeholder">Tady budou vidět přidaná pole</i>`)
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
      case "autoři":
        fieldInitFunc = initAuthorsField;
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
  function initAuthorsField(field) {
    // Create the HTML representing the element
    field.innerHTML = `
      <span class="f-start">Autoři:</span>
      <div class="flex-column max-flex f-middle authors-field">
        <span class="flex-row case-controls">
          <b>Písmo jména (první autor):</b>
          <label class="flex-row"><input type="radio" name="a-formatjmeno-prvni" value="vychozi">Výchozí</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno-prvni" value="male">Malé</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno-prvni" value="velke">Velké</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno-prvni" value="prvnivelke" checked>První&nbsp;velké</label>
        </span>
        <span class="flex-row case-controls">
          <b>Písmo příjmení (první autor):</b>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni-prvni" value="vychozi">Výchozí</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni-prvni" value="male">Malé</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni-prvni" value="velke" checked>Velké</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni-prvni" value="prvnivelke">První&nbsp;velké</label>
        </span>
        <span class="flex-row case-controls">
          <b>Na prvním místě (první autor):</b>
          <label class="flex-row"><input type="radio" name="a-poradi-prvni" value="prijmeni" checked>Příjmení</label>
          <label class="flex-row"><input type="radio" name="a-poradi-prvni" value="jmeno">Jméno</label>
        </span>
        <label class="flex-row"><b>Interpunkce&nbsp;mezi&nbsp;jmény (první autor):</b><input class="max-flex" type="text" value="," name="a-intjmeno-prvni"></label>
        <hr class="max-flex">
        <span class="flex-row case-controls">
          <b>Písmo jména:</b>
          <label class="flex-row"><input type="radio" name="a-formatjmeno" value="vychozi">Výchozí</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno" value="male">Malé</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno" value="velke">Velké</label>
          <label class="flex-row"><input type="radio" name="a-formatjmeno" value="prvnivelke" checked>První&nbsp;velké</label>
        </span>
        <span class="flex-row case-controls">
          <b>Písmo příjmení:</b>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni" value="vychozi">Výchozí</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni" value="male">Malé</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni" value="velke" checked>Velké</label>
          <label class="flex-row"><input type="radio" name="a-formatprijmeni" value="prvnivelke">První&nbsp;velké</label>
        </span>
        <span class="flex-row case-controls">
          <b>Na prvním místě (ostatní):</b>
          <label class="flex-row"><input type="radio" name="a-poradi" value="prijmeni" checked>Příjmení</label>
          <label class="flex-row"><input type="radio" name="a-poradi" value="jmeno">Jméno</label>
        </span>
        <label class="flex-row"><b>Interpunkce&nbsp;mezi&nbsp;jmény (ostatní):</b><input class="max-flex" type="text" value="," name="a-intjmeno"></label>
        <hr class="max-flex">
        <label class="flex-row"><b>Interpunkce&nbsp;mezi&nbsp;autory:</b><input class="max-flex" type="text" value=";" name="a-intautor"></label>
        <span class="flex-row case-controls">
          <b>Spojka před posledním autorem:</b>
          <label class="flex-row"><input type="radio" name="a-a" value="a" checked>a</label>
          <label class="flex-row"><input type="radio" name="a-a" value="&amp;">&amp;</label>
          <label class="flex-row"><input type="radio" name="a-a" value="and">and</label>
        </span>
        <span class="flex-row case-controls">
          <b>Použít spojku vždy když je více než jeden autor:</b>
          <label class="flex-row"><input type="radio" name="a-vzdya" value="" checked>Ne</label>
          <label class="flex-row"><input type="radio" name="a-vzdya" value="1">Ano</label>
        </span>
        <label class="flex-row"><b>Maximalní&nbsp;počet&nbsp;autorů:</b><input class="max-flex" type="number" min="1" step="1" value="5" name="a-max"></label>
        <label class="flex-row"><b>Přípona (a další):</b><input class="max-flex" type="text" value="et al." name="a-etal"></label>
        <hr class="max-flex">
        <div class="flex-row max-flex">
          ${fieldSeparatorFormControls}
          ${addSpaceFormControls}
        </dev>
      </div>
    `;

    // Stop text input fields from being dragable
    stopElementFromBeingDragged(field.elements.namedItem("f-oddělovač"), field);
    stopElementFromBeingDragged(
      field.elements.namedItem("a-intjmeno-prvni"),
      field
    );
    stopElementFromBeingDragged(field.elements.namedItem("a-intjmeno"), field);
    stopElementFromBeingDragged(field.elements.namedItem("a-intautor"), field);
    stopElementFromBeingDragged(field.elements.namedItem("a-etal"), field);

    // Put the function necessary for rendering the part of template
    // represented by the field into global weakMap so it can be called
    // from the citation generator part of this sphagetti
    const data = fieldCustomData.get(field);
    data.getTemplateValue = function () {
      /** @type {string[]} Arguments for the autoři helper */
      const args = [];
      const inputs = field.elements;

      for (const input of inputs) {
        if ("value" in input && "name" in input) {
          /** @type {string} */
          let key = input.name;
          if (key.startsWith("a-")) {
            key = key.slice(2); // Remove the "a-" prefix
          } else if (key.startsWith("f-")) {
            // Formating control, ignore
            continue;
          } else {
            // If you see this warning then check the HTML template above.
            console.warn(
              "initAuthorsField found input with name without 'a-' prefix. Got:",
              key
            );
            continue;
          }
          // Skip unchecked radio inputs
          if (
            "type" in input &&
            input.type === "radio" &&
            "checked" in input &&
            !input.checked
          ) {
            continue;
          }
          console.log(input, input.value);
          args.push(`${key}="${input.value}"`);
        }
      }

      let expr = "{{autoři " + args.join(" ") + "}}";
      expr = addSeparator(expr, field);
      expr = addSpace(expr, field);
      return expr;
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
    initGenericTimeField(field, "Datum&nbsp;vydání", "datum-vydání");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initUrlField(field) {
    initGenericUrlField(field, "URL", "url");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initArchivalUrlField(field) {
    initGenericUrlField(field, "Archivní&nbsp;URL", "archivní-url");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initDateOfHarvestField(field) {
    initGenericTimeField(field, "Datum archivace", "datum-archivace");
  }

  /**
   * @param {HTMLFormElement} field
   */
  function initDatefCitationField(field) {
    initGenericTimeField(field, "Datum citace", "datum-citace");
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
            <b>Písmo:</b>
          ${fieldFormatFormControls}
          ${fieldCaseFormControls}
        </div>
        <div class="flex-row">
          ${fieldSeparatorFormControls}
          ${addSpaceFormControls}
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
      expr = addSpace(expr, field);
      return expr;
    };
  }

  /**
   * @param {HTMLFormElement} field
   * @param {string} readableName
   * @param {string} exprName
   */
  function initGenericTimeField(field, readableName, exprName) {
    field.innerHTML = `
      <span class="f-start">${readableName}:</span>
      <div class="flex-column max-flex f-middle">
        <div class="flex-row max-flex">
          <select class="flex-row" name="f-date-format">
            <option value="iso-date">Datum (RRRR-MM-DD)</option>
            <option value="rok">Rok</option>
            <option value="human">Datum dlouhé</option>
            <option value="iso">Datum a čas (ISO 8601)</option>
            <option value="rfc">Datum a čas (ISO 8601 s mezerou)</option>
            <option value="apa">APA (RRRR, měsíc DD)</option>
            <option value="bez-formatu">Neměnit formát</option>
          </select>
          <label class="flex-row"><input type="checkbox" name="f-utc">UTC</label>
        </div>
        <div class="flex-row max-flex">
          ${addSpaceFormControls}
        </div>
      </div>
    `;
    //<option value="iso-time">ISO 8601 jen čas</option>
    const data = fieldCustomData.get(field);
    data.getTemplateValue = function () {
      let expr = exprName;
      expr = formatTime(field, expr);
      expr = `{{${expr}}}`;
      expr = addSpace(expr, field);
      return expr;
    };

    /**
     * @param {HTMLFormElement} field
     * @param {string} expr
     */
    function formatTime(field, expr) {
      expr = `datum ${expr}`;
      let formatParam = "";
      const utc = field.elements.namedItem("f-utc");
      if (utc !== null && "checked" in utc && utc.checked) {
        formatParam = "utc-";
      }
      const format = field.elements.namedItem("f-date-format");
      if (format !== null && "value" in format) {
        formatParam += format.value;
        expr += ' "' + formatParam + '"';
      }
      return expr;
    }
  }

  /**
   * @param {HTMLFormElement} field
   * @param {string} readableName
   * @param {string} exprName
   */
  function initGenericUrlField(field, readableName, exprName) {
    field.innerHTML = `
      <span class="f-start">${readableName}:</span>
      <div class="flex-column max-flex f-middle">
        <div class="flex-row max-flex">
          ${fieldFormatFormControls}
          ${fieldSeparatorFormControls}
          ${addSpaceFormControls}
        </div>
      </div>
    `;
    stopElementFromBeingDragged(field.elements.namedItem("f-oddělovač"), field);
    const data = fieldCustomData.get(field);
    data.getTemplateValue = function () {
      let expr = exprName;
      expr = wrapExprInFormat(expr, field);
      expr = `{{${expr}}}`;
      expr = addSeparator(expr, field);
      expr = addSpace(expr, field);
      return expr;
    };
  }

  const addSpaceFormControls = `<label class="flex-row"><input type="checkbox" name="f-add-space" checked>Přidat&nbsp;mezeru</label>`;
  /**
   * Will append space at the end of expr if f-add-space is checked.
   * @param {string} expr
   * @param {HTMLFormElement} field
   * @return {string}
   */
  function addSpace(expr, field) {
    const addSpaceInput = field.elements.namedItem("f-add-space");
    if (
      addSpaceInput !== null &&
      "checked" in addSpaceInput &&
      addSpaceInput.checked
    ) {
      return expr + " ";
    }
    return expr;
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
    <label class="flex-row"><b>Interpunkce:</b><input type="text" name="f-oddělovač"></label>
  `;
  /**
   * Will append separator at the end of expr.
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
    <label class="flex-row"><input type="radio" name="f-case" value="no-change" checked>Výchozí</label>
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
        helperName = "velké";
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
   * Create the template string from fields
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
      template += data.getTemplateValue();
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

    // Part of the form contatining authors metadata
    const authorsDiv = document.getElementById("authors");
    if (authorsDiv === null) {
      throw new Error("Element with id 'authors' must exist");
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
          fillForm(generatorForm, authorsDiv, inputData[0]);
          enableFormControls(
            generatorForm,
            authorsDiv,
            templateElement,
            citationOutput,
            inputData
          );
        }
      } catch (err) {
        console.error(err);
      }
    }

    // Copy values from datepickers to their text input elements
    enableDatepickers();

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
   * @param {HTMLElement} authorsDiv
   * @param {HTMLElement} templateElement
   * @param {HTMLElement} citationOutput
   * @param {Array<any>} citationData
   */
  function enableFormControls(
    generatorForm,
    authorsDiv,
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
      fillForm(generatorForm, authorsDiv, citationData[currentDataIndex]);
      generateCitation(generatorForm, templateElement, citationOutput);
    });

    nextBtn.addEventListener("click", () => {
      if (currentDataIndex >= dataCount - 1) {
        return;
      }
      saveCurrentCitationData(generatorForm, citationData, currentDataIndex);
      currentDataIndex++;
      currentIndexElem.textContent = (currentDataIndex + 1).toString();
      fillForm(generatorForm, authorsDiv, citationData[currentDataIndex]);
      generateCitation(generatorForm, templateElement, citationOutput);
    });

    const addAuthorButton = document.getElementById("add-author");
    if (addAuthorButton === null) {
      throw new Error("Element with id 'add-author' must exist");
    }
    addAuthorButton.addEventListener("click", () => {
      addAuthorFieldset(authorsDiv);
      generatorForm.dispatchEvent(new Event("input", { bubbles: true }));
    });

    const removeAuthorButton = document.getElementById("remove-author");
    if (removeAuthorButton === null) {
      throw new Error("Element with id 'remove-author' must exist");
    }
    removeAuthorButton.addEventListener("click", () => {
      removeLastAuthorFieldset(authorsDiv);
      generatorForm.dispatchEvent(new Event("input", { bubbles: true }));
    });
  }

  /**
   * Register input callbacks to the datetime input elements so they propagate
   * the value to the adjacent text input elements.
   * This is done so that user has the power to use custom date format but also
   * has the convinience of datepicker.
   */
  function enableDatepickers() {
    const elementIDs = [
      { textInput: "datum-vydání", dateTimeInput: "datum-vydání-datetime" },
      {
        textInput: "datum-archivace",
        dateTimeInput: "datum-archivace-datetime",
      },
      { textInput: "datum-citace", dateTimeInput: "datum-citace-datetime" },
    ];
    elementIDs.forEach((ids) => {
      const textInput = document.getElementById(ids.textInput);
      if (textInput === null) {
        throw new Error(`Element with id ${ids.textInput} must exist`);
      }
      const dateTimeInput = document.getElementById(ids.dateTimeInput);
      if (dateTimeInput === null) {
        throw new Error(`Element with id ${ids.dateTimeInput} must exist`);
      }

      dateTimeInput.addEventListener("change", () => {
        textInput.value = dateTimeInput.value;
      });
    });
  }

  /**
   * Save current form values into citationData array
   * @param {HTMLFormElement} generatorForm
   * @param {Array<any>} citationData
   * @param {number} currentIndex
   */
  function saveCurrentCitationData(generatorForm, citationData, currentIndex) {
    citationData[currentIndex] = extractCurrentCitationData(generatorForm);
  }

  /**
   * Extract data from generatorForm and return them
   * @param {HTMLFormElement} generatorForm
   * @return {object}
   */
  function extractCurrentCitationData(generatorForm) {
    const currentData = {};
    currentData.autoři = [];
    for (const element of generatorForm.elements) {
      // Handle all fieldsets with author data
      if (
        "authorId" in element.dataset &&
        element instanceof HTMLFieldSetElement
      ) {
        const author = {};
        const firstname = element.elements.namedItem("jméno")?.value;
        const lastname = element.elements.namedItem("příjmení")?.value;
        author.jméno = firstname ?? "";
        author.příjmení = lastname ?? "";
        // This uses the fact that js will not error on inserting to random array indices.
        // It is not the responsibility of this function to ensure that the author fieldsets are properly numbered and ordered.
        const authorId = Number(element.dataset.authorId);
        currentData.autoři[authorId] = author;
        // Hanlde the rest of data
      } else if ("citationfield" in element.dataset) {
        currentData[element.dataset.citationfield] = element.value;
      }
    }
    return currentData;
  }

  /**
   *
   * @param {HTMLFormElement} generatorForm
   * @param {HTMLElement} authorsDiv // Part of form for author fieldsets
   * @param {any} inputData // Object containing citation field values
   */
  function fillForm(generatorForm, authorsDiv, inputData) {
    authorsDiv.innerHTML = "";
    for (const field in inputData) {
      if (field === "autoři") {
        inputData[field].forEach((author, id) => {
          const authorHTML = createAuthorHTML(author, id);
          authorsDiv.insertAdjacentHTML("beforeend", authorHTML);
        });
      } else {
        const element = generatorForm.elements.namedItem(field);
        if (element === null) {
          console.warn(`Form element with id ${field} does not exist!`);
          continue;
        }
        element.value = inputData[field];
      }
    }
  }

  /**
   * @param {{jméno: string, příjmení: string}} author
   * @param {number} id
   * @returns {string} String containing HTML representing the author form element
   */
  function createAuthorHTML(author, id) {
    return `
      <fieldset data-author-id="${id}">
        <legend>Autor ${id + 1}</legend>
        <div class="flex-row cit-gen-fields">
          <div class="flex-column cit-gen-labels">
            <label for="příjmení">Příjmení:</label>
            <label for="jméno">Jméno:</label>
          </div>
          <div class="flex-column cit-gen-inputs">
            <input type="text" name="příjmení" value=${author.příjmení}>
            <input type="text" name="jméno" value=${author.jméno}>
          </div>
        </div>
      </fieldset>
    `;
  }

  /**
   * Add empty author fieldset to the form
   * @param {HTMLElement} authorsDiv
   */
  function addAuthorFieldset(authorsDiv) {
    let id = authorsDiv.lastElementChild?.dataset.authorId;
    if (id) {
      id = Number(id) + 1;
    } else {
      id = 0;
    }
    authorsDiv.insertAdjacentHTML(
      "beforeend",
      createAuthorHTML({ jméno: "", příjmení: "" }, id)
    );
  }

  /**
   * Remove last author fieldset from the form
   * @param {HTMLElement} authorsDiv
   */
  function removeLastAuthorFieldset(authorsDiv) {
    const lastElement = authorsDiv.lastElementChild;
    if (lastElement) {
      lastElement.remove();
    }
    // Otherwise do nothing
  }

  /**
   * Replace citationOutput element with template and data from generatorForm.
   * @param {HTMLFormElement} generatorForm
   * @param {HTMLElement} templateElement
   * @param {HTMLElement} citationOutput
   */
  function generateCitation(generatorForm, templateElement, citationOutput) {
    const template = Handlebars.compile(templateElement.value);
    const data = extractCurrentCitationData(generatorForm);
    const citation = template(data);
    citationOutput.innerHTML = citation;
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
  Handlebars.registerHelper("velké", upperCase);
  Handlebars.registerHelper("první-velké", capitalize);
  Handlebars.registerHelper("malé", lowerCase);
  // Take array of authors and format them to ISO 690
  Handlebars.registerHelper("autoři", formatAuthors);
  Handlebars.registerHelper("datum", formatDateHelper);

  /**
   * @param {boolean} end
   * @param {any} text
   */
  function formatBold(end, text) {
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
    if (text.length === 0) {
      return text;
    }
    if (text.length === 1) {
      return text.toUpperCase();
    }
    return text
      .toLowerCase()
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

  /**
   * This helper will fill the template with authors formated per ISO and user arguments.
   * @param  {any} hashOptions hashOptions.hash contains object with options
   */
  function formatAuthors(hashOptions) {
    /*
    Output can look like (iso - other styles are very similar and usually simpler):
    LASTNAME, Firstname; LASTNAME, Firstname a LASTNAME, Firstname.
    LASTNAME, Firstname; LASTNAME, Firstname; ... LASTNAME, Firstname et al.
    */

    /*
    Posible options:
    number of authors to show
    should we show "et al." or something else when there is more
    firsname lastname separators
    author separator
    firsname format - case, bold/italic
    lastname format
    different format for fisrt author
    add "a" or "&" before last author
    order of firstname and lastname
    */

    /** @type {{jméno:string, příjmení:string}[]} List containing author objects from template context */
    const authors = this.autoři ?? [];
    /** @type {string[]} List of strings that will be joined to produce final output */
    const outputBuilder = [];
    /** @type {object} Unwrapped options */
    const options = hashOptions.hash;

    /** How many authors to show (five is default as in ISO 690) */
    const maxAuthors = options.max !== undefined ? Number(options.max) : 5;
    /** @type {string} What to show if we have more than max authors */
    const overLimitSuffix =
      options.etal !== undefined ? options.etal : "et al.";
    /** @type {string} Separator between lastname and firstname (first author) */
    const nameSeparatorFirstAuthor =
      options["intjmeno-prvni"] !== undefined ? options["intjmeno-prvni"] : ",";
    /** @type {string} Separator between lastname and firstname */
    const nameSeparator =
      options.intjmeno !== undefined ? options.intjmeno : ",";
    /** @type {string} Separator between authors */
    const authorSeparator =
      options.intautor !== undefined ? options.intautor : ";";
    /** @type {string} Firstname formatting options */
    const firstnameFormatFirstAuthor =
      options["formatjmeno-prvni"] !== undefined
        ? options["formatjmeno-prvni"]
        : "prvnivelke";
    /** @type {string} Firstname formatting options */
    const firstnameFormatOther =
      options.formatjmeno !== undefined ? options.formatjmeno : "prvnivelke";
    /** @type {string} Lastname formatting options */
    const lastnameFormatFirstAuthor =
      options["formatprijmeni-prvni"] !== undefined
        ? options["formatprijmeni-prvni"]
        : "velke";
    /** @type {string} Lastname formatting options */
    const lastnameFormatOther =
      options.formatprijmeni !== undefined ? options.formatprijmeni : "velke";
    /** @type {string} What to add before last author (ussualy "a" or "&" or "and") */
    const andSeparator = options.a !== undefined ? options.a : "a";
    /** @type {boolean} Should we print andSeparator even if we have max authors. */
    const alwaysPrintAndSeparator =
      options.vzdya !== undefined ? Boolean(options.vzdya) : false;
    /** @type {string} Order of firsname and lastname (first author) */
    const nameOrderFirtsAuthor =
      options["poradi-prvni"] !== undefined
        ? options["poradi-prvni"]
        : "prijmeni";
    /** @type {string} Order of firsname and lastname (rest) */
    const nameOrder =
      options.poradi !== undefined ? options.poradi : "prijmeni";

    for (let i = 0; i < authors.length; i++) {
      // This will be the smaller of the two possible limits on number of authors
      const authorLimit = Math.min(authors.length, maxAuthors);

      // If author limit break loop
      if (i >= authorLimit) {
        // If overLimitSufix is set and maxAuthors is smaller than authors.length
        if (overLimitSuffix !== "" && maxAuthors < authors.length) {
          outputBuilder.push(" ", overLimitSuffix);
        }
        break;
      }

      const author = authors[i];
      // Ensure that name parts exist
      author.jméno ??= "";
      author.příjmení ??= "";

      let authorStr = "";

      // Select currect format to use
      let firstnameFormat = firstnameFormatOther;
      let lastnameFormat = lastnameFormatOther;
      if (i === 0) {
        firstnameFormat = firstnameFormatFirstAuthor;
        lastnameFormat = lastnameFormatFirstAuthor;
      }

      // Set lettercase of name parts
      const firstname = chooseCaseFunc(firstnameFormat)(author.jméno);
      const lastname = chooseCaseFunc(lastnameFormat)(author.příjmení);

      // This is used to check if the name part was empty
      let lastAddedNamePart = "";

      // Get the correct value to determine the order of firstname and lastname
      let nameOrderThisIter = nameOrder;
      if (i === 0) {
        nameOrderThisIter = nameOrderFirtsAuthor;
      }

      // If order is set to "prijmeni", then lastname will go first
      if (nameOrderThisIter === "prijmeni") {
        authorStr += lastname;
        lastAddedNamePart = lastname;
      } else {
        authorStr += firstname;
        lastAddedNamePart = firstname;
      }

      // Set correct name separator
      let nameSeparatorThisIter = nameSeparator;
      if (i === 0) {
        nameSeparatorThisIter = nameSeparatorFirstAuthor;
      }

      // Add name separator and space (this space should be added always)
      if (lastAddedNamePart !== "") {
        authorStr += nameSeparatorThisIter + " ";
      } else {
        authorStr += " "; // No name part, add only space
      }

      // Add rest of the name
      if (nameOrderThisIter === "prijmeni") {
        authorStr += firstname;
        lastAddedNamePart = firstname;
      } else {
        authorStr += lastname;
        lastAddedNamePart = lastname;
      }

      // If andSeparator is set
      // and there is at least two authors
      // and we are before the last author
      // and the number of authors to print is less then authors limit or the alwaysPrintAndSeparator is true
      // then print the andSeparator
      if (
        andSeparator !== "" &&
        authors.length > 1 &&
        i === authorLimit - 2 &&
        (authors.length <= authorLimit || alwaysPrintAndSeparator)
      ) {
        authorStr += " " + andSeparator + " ";
      } else {
        // Otherwise try to print author separator

        // If we are not last
        // (and from before we know that the conditions for andSeparator are not met)
        // then print authorSeparator
        if (i < authorLimit - 1) {
          if (lastAddedNamePart !== "") {
            authorStr += authorSeparator + " ";
          } else {
            authorStr += " ";
          }
        }
      }

      // Add to builder
      outputBuilder.push(authorStr);
    }

    return outputBuilder.join("");
  }

  /**
   * Used for author formatting
   * @param {string} formatOptions
   * @returns {function(string):string} Function that takes text transforms letter case and returns it
   */
  function chooseCaseFunc(formatOptions) {
    const options = new Set(formatOptions.split(","));
    if (options.has("velke")) {
      return upperCase;
    } else if (options.has("prvnivelke")) {
      return capitalize;
    } else if (options.has("male")) {
      return lowerCase;
    } else {
      return function (text) {
        return text;
      };
    }
  }

  /**
   * Helper for formating dates and times.
   * @param {string} date Date and time in simplified iso format or custom user data.
   * @param {string} [format] Name of format to use. If not specified, then date is returned unchanged. The string may begin with "utc-", then the output will be formated with utc timezone.
   */
  function formatDateHelper(date, format) {
    date = Handlebars.escapeExpression(date);
    if (format === undefined) {
      return date;
    }

    let parsedDate = luxon.DateTime.fromISO(date);

    if (format.startsWith("utc-")) {
      format = format.slice(4);
      parsedDate = parsedDate.toUTC();
    }

    switch (format) {
      case "iso":
        return parsedDate.toISO();
      case "iso-date":
        return parsedDate.toISODate();
      case "iso-time":
        return parsedDate.toISOTime();
      case "rfc":
        return parsedDate.toFormat("yyyy-MM-dd HH:mm:ssZZ");
      case "rok":
        return parsedDate.toFormat("yyyy");
      case "apa":
        return parsedDate.toFormat("yyyy, MMMM d");
      case "human":
        return parsedDate.toLocaleString(luxon.DateTime.DATE_FULL);
      default:
        console.error(
          `formatDateHelper recieved invalid format argument: ${format}`
        );
        return date; // Return date so user sees something.
    }
  }

  main();
})();
