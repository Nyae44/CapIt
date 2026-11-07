import { saveCaption } from "./storage.js";

// DOM Elements
const form = document.querySelector("#caption-form");
const descriptionInput = document.querySelector("#description");
const descriptionField = descriptionInput.closest(".field");
const descriptionError = document.querySelector("#description-error");

const generateBtn = document.querySelector("#generate-btn");
const emptyState = document.querySelector("#empty-state");
const result = document.querySelector("#result");
const captionText = document.querySelector("#caption-text");

const copyBtn = document.querySelector("#copy-btn");
const editBtn = document.querySelector("#edit-btn");
const regenerateBtn = document.querySelector("#regenerate-btn");

const alternativesList = document.querySelector("#alternatives-list");
const alternativesBox = document.querySelector("#alternatives");
const toast = document.querySelector("#toast");

// API Configuration
const API_URL = "/api/v1/captions/generate";

const GENERATE_LABEL = generateBtn.textContent;
const REGENERATE_LABEL = regenerateBtn.textContent;

// Application State
let currentOptions = null;
let currentCaptions = [];

// Read form options
function readOptions() {
    const data = new FormData(form);

    return {
        description: String(data.get("description") ?? ""),
        platform: data.get("platform"),
        contentType: data.get("contentType"),
        tone: data.get("tone"),
        language: data.get("language"),
        length: data.get("length"),
        instructions: String(data.get("instructions") ?? "")
    };
}

// Convert frontend options into backend API format
function buildPayload(options) {
    return {
        platform: options.platform,
        content_type: options.contentType,
        topic: options.description,
        tone: options.tone,
        language: options.language,
        caption_length: options.length
    };
}

// API Integration
async function generateCaptions(options) {

    const payload = buildPayload(options);

    const response = await fetch(API_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
        throw new Error(
            data.error ||
            data.message ||
            `Caption generation failed (${response.status})`
        );
    }

    // Support either a single caption or multiple captions.
    let captions = [];

    if (Array.isArray(data.captions)) {
        captions = data.captions;
    } else if (typeof data.caption === "string") {
        captions = [data.caption];
    } else if (typeof data.content === "string") {
        captions = [data.content];
    } else if (typeof data === "string") {
        captions = [data];
    }

    captions = captions.filter(
        caption => typeof caption === "string" && caption.trim()
    );

    if (captions.length === 0) {
        throw new Error("The server returned no captions.");
    }

    return captions;
}

// Loading State
function setLoading(button, isLoading, idleLabel) {
    button.disabled = isLoading;

    button.setAttribute(
        "aria-busy",
        String(isLoading)
    );

    button.innerHTML = isLoading
        ? '<span class="spinner" aria-hidden="true"></span>Writing...'
        : idleLabel;
}

// Toast Notifications
function showToast(message) {
    toast.textContent = message;
    toast.classList.add("is-visible");

    clearTimeout(showToast.timer);

    showToast.timer = setTimeout(() => {
        toast.classList.remove("is-visible");
    }, 2000);
}

// Render Alternative Captions
function renderAlternatives() {
    alternativesList.replaceChildren();

    const current = captionText.value;

    const alternatives = currentCaptions
        .filter(caption => caption !== current)
        .slice(0, 3);

    alternativesBox.hidden = alternatives.length === 0;

    for (const caption of alternatives) {
        const item = document.createElement("li");

        const button = document.createElement("button");

        button.type = "button";
        button.className = "alternative-btn";
        button.textContent = caption;

        button.addEventListener("click", () => {
            captionText.value = caption;

            renderAlternatives();

            showToast("Alternative applied");
        });

        item.append(button);
        alternativesList.append(item);
    }
}

// Display Generated Caption
function showResult(captions) {
    currentCaptions = captions;

    captionText.value = captions[0];

    captionText.readOnly = true;

    editBtn.textContent = "Edit";
    editBtn.setAttribute("aria-pressed", "false");

    emptyState.hidden = true;
    result.hidden = false;

    renderAlternatives();
}

// Save Caption
function saveCurrentCaption(caption) {
    saveCaption({
        content: caption,
        platform: currentOptions.platform,
        tone: currentOptions.tone
    });
}

// Generate Caption
async function runGeneration(button, idleLabel) {

    if (!currentOptions) {
        return;
    }

    setLoading(button, true, idleLabel);

    try {

        const captions = await generateCaptions(currentOptions);

        showResult(captions);

        saveCurrentCaption(captions[0]);

        showToast("Caption generated successfully!");

    } catch (error) {

        console.error("Caption generation error:", error);

        showToast(error.message || "Something went wrong. Please try again.");

    } finally {

        setLoading(button, false, idleLabel);

    }
}

// Form Submission
form.addEventListener("submit", (event) => {

    event.preventDefault();

    // Validate description
    if (!descriptionInput.value.trim()) {

        descriptionField.classList.add("has-error");

        descriptionError.hidden = false;

        descriptionInput.setAttribute("aria-invalid", "true");

        descriptionInput.focus();

        return;
    }

    descriptionField.classList.remove("has-error");

    descriptionError.hidden = true;

    descriptionInput.removeAttribute("aria-invalid");

    // Collect form data
    currentOptions = readOptions();

    // Send request to Go backend
    runGeneration(generateBtn, GENERATE_LABEL);

});

// Description Validation
descriptionInput.addEventListener("input", () => {

    if (descriptionInput.value.trim()) {

        descriptionField.classList.remove("has-error");

        descriptionError.hidden = true;

        descriptionInput.removeAttribute("aria-invalid");

    }

});

// Regenerate Caption
regenerateBtn.addEventListener("click", () => {

    if (currentOptions) {
        runGeneration(regenerateBtn, REGENERATE_LABEL);
    }

});

// Edit Caption
editBtn.addEventListener("click", () => {

    const isEditing = captionText.readOnly;

    captionText.readOnly = !isEditing;

    editBtn.textContent = isEditing ? "Done" : "Edit";

    editBtn.setAttribute(
        "aria-pressed",
        String(isEditing)
    );

    if (isEditing) {

        captionText.focus();

        captionText.setSelectionRange(
            captionText.value.length,
            captionText.value.length
        );

    } else {

        renderAlternatives();

    }

});

// Copy Caption
copyBtn.addEventListener("click", async () => {

    try {

        await navigator.clipboard.writeText(captionText.value);

        showToast("Caption copied to clipboard");

    } catch (error) {

        console.error("Clipboard error:", error);

        showToast("Couldn't copy. Select the text and copy it manually.");

    }

});