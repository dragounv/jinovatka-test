// The generator view script
(function () {
  function main() {
    prepareTemplateBuilder();
    prepareCitationGenerator();
  }

  function prepareTemplateBuilder() {}

  function prepareCitationGenerator() {
    // Get elements for input form and output paragraph.
    const generatorForm = document.getElementById("generator");
    const citationOutput = document.getElementById("citation");
    const inputDataElement = document.getElementById("input-data");

    if (!(generatorForm instanceof HTMLFormElement)) {
      throw new TypeError(
        "Element with id 'generator' must be HTMLFormElement"
      );
    }

    if (citationOutput === null) {
      throw new TypeError("Element with id 'citation' must exist");
    }

    // If input data are present parse them and fill the Form.
    if (inputDataElement !== null) {
      try {
        if (inputDataElement.type !== "application/json") {
          throw new TypeError(
            "Element #input-data must be script with type 'application/json'"
          );
        }
        const inputData = JSON.parse(inputDataElement.text);
        if (!Array.isArray(inputData)) {
          throw new TypeError("inputData must be array of objects");
        }
        if (inputData.length !== 0) {
          fillForm(generatorForm, inputData[0]);
          enableFormControls(generatorForm, citationOutput, inputData);
        }
      } catch (err) {
        console.error(err);
      }
    }

    // Render the template any time when user inputs data.
    generatorForm.addEventListener("input", () =>
      generateCitation(generatorForm, citationOutput)
    );

    // Render the template first time on page load.
    generateCitation(generatorForm, citationOutput);
  }

  /**
   *
   * @param {HTMLFormElement} generatorForm
   * @param {HTMLElement} citationOutput
   * @param {Array} citationData
   */
  function enableFormControls(generatorForm, citationOutput, citationData) {
    const formControls = document.getElementById("form-controls");
    if (formControls === null) {
      throw new TypeError("Element with id 'form-controls' must exist");
    }

    formControls.hidden = false;
    formControls.classList.remove("hidden");

    let currentDataIndex = 0;
    let dataCount = citationData.length;

    const currentIndexElem = document.getElementById("cit-data-num");
    const countElem = document.getElementById("cit-data-count");
    if (currentIndexElem === null) {
      throw new TypeError("Element with id 'cit-data-num' must exist");
    }
    if (countElem === null) {
      throw new TypeError("Element with id 'cit-data-count' must exist");
    }

    // Show index bigger by one as that is what people generaly expect
    currentIndexElem.textContent = (currentDataIndex + 1).toString();
    countElem.textContent = dataCount.toString();

    const prevBtn = document.getElementById("prev");
    const nextBtn = document.getElementById("next");
    if (prevBtn === null) {
      throw new TypeError("Element with id 'prev' must exist");
    }
    if (nextBtn === null) {
      throw new TypeError("Element with id 'next' must exist");
    }

    prevBtn.addEventListener("click", () => {
      if (currentDataIndex <= 0) {
        return;
      }
      saveCurrentCitationData(generatorForm, citationData, currentDataIndex);
      currentDataIndex--;
      currentIndexElem.textContent = (currentDataIndex + 1).toString();
      fillForm(generatorForm, citationData[currentDataIndex]);
      generateCitation(generatorForm, citationOutput);
    });

    nextBtn.addEventListener("click", () => {
      if (currentDataIndex >= dataCount - 1) {
        return;
      }
      saveCurrentCitationData(generatorForm, citationData, currentDataIndex);
      currentDataIndex++;
      currentIndexElem.textContent = (currentDataIndex + 1).toString();
      fillForm(generatorForm, citationData[currentDataIndex]);
      generateCitation(generatorForm, citationOutput);
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
   * @param {HTMLElement} citationOutput
   */
  function generateCitation(generatorForm, citationOutput) {
    const template = Handlebars.compile(
      generatorForm.elements["template"].value
    );
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

  main();
})();
