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
          fillForm(generatorForm, inputData);
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
   * @param {Array} inputData
   */
  function fillForm(generatorForm, inputData) {
    const firstCitationData = inputData[0];
    for (const field in firstCitationData) {
      const element = generatorForm.elements.namedItem(field);
      if (element === null) {
        console.warn(`Form element with id ${field} does not exist!`);
        continue;
      }
      element.value = firstCitationData[field];
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
