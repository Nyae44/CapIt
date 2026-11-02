import { saveCaption } from "./storage.js";
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

const GENERATE_LABEL = generateBtn.textContent;
const REGENERATE_LABEL = regenerateBtn.textContent;

const TONE_OPENERS = {
  friendly: ["Here's a little something to brighten your feed.", "Come along with us!", "We're so happy to share this with you."],
  professional: ["Introducing a closer look at our latest work.", "A considered update worth your attention.", "Here is what we have been building."],
  playful: ["Okay, this one made us smile!", "Plot twist: it's even better in person.", "Warning: may cause happy scrolling."],
  inspirational: ["Every big story starts with one small step.", "Dream it. Share it. Make it happen.", "Today is a good day to begin again."],
  witty: ["Nobody asked, but here we are.", "We planned this. Mostly.", "Our favorite kind of plan: the spontaneous one."],
  bold: ["Stop scrolling. This matters.", "No hedging. Just results.", "Bold moves only."],
};

const CTAS = {
  instagram: "Double tap if you agree and tell us in the comments.",
  tiktok: "Follow for more and drop a comment below.",
  linkedin: "What would you add? Share your thoughts below.",
  x: "Reply with your take.",
  facebook: "Tag a friend who would love this.",
};

const HASHTAGS = {
  instagram: ["#dailyinspo", "#contentcreator", "#socialmedia"],
  tiktok: ["#fyp", "#trending", "#creatorsofinstagram"],
  linkedin: ["#growth", "#leadership", "#insights"],
  x: ["#thoughts", "#today"],
  facebook: ["#community", "#shareit"],
};

const LENGTH_HASHTAG_COUNT = { short: 1, medium: 2, long: 3 };

function firstSentence(text) {
  const trimmed = text.trim().replace(/\s+/g, " ");
  const cleaned = trimmed.replace(/[.!?]+$/, "");
  return cleaned.charAt(0).toUpperCase() + cleaned.slice(1);
}

function buildCaption(options, variant) {
  const openers = TONE_OPENERS[options.tone];
  const opener = openers[variant % openers.length];
  const subject = firstSentence(options.description);
  const hashtags = HASHTAGS[options.platform].slice(0, LENGTH_HASHTAG_COUNT[options.length]).join(" ");
  const cta = CTAS[options.platform];
  const extra = options.instructions.trim();

  const lines = [opener, `${subject}.`];
  if (options.length !== "short") {
    lines.push(extra ? `${extra.replace(/[.!?]+$/, "")}.` : "Every moment, shared with intention.");
  }
  if (options.length === "long") {
    lines.push("We hope it sparks something for you too.");
  }
  lines.push(cta);
  if (hashtags) lines.push(hashtags);
  return lines.join("\n\n");
}

function generateCaptions(options) {
  return new Promise((resolve) => {
    setTimeout(() => {
      const captions = [0, 1, 2, 3].map((variant) => buildCaption(options, variant));
      resolve(captions);
    }, 600);
  });
}

function readOptions() {
  const data = new FormData(form);
  return {
    description: String(data.get("description") ?? ""),
    platform: data.get("platform"),
    contentType: data.get("contentType"),
    tone: data.get("tone"),
    language: data.get("language"),
    length: data.get("length"),
    instructions: String(data.get("instructions") ?? ""),
  };
}

let currentOptions = null;
let currentCaptions = [];

function setLoading(button, isLoading, idleLabel) {
  button.disabled = isLoading;
  button.setAttribute("aria-busy", String(isLoading));
  button.innerHTML = isLoading ? '<span class="spinner" aria-hidden="true"></span>Writing...' : idleLabel;
}

function showToast(message) {
  toast.textContent = message;
  toast.classList.add("is-visible");
  clearTimeout(showToast.timer);
  showToast.timer = setTimeout(() => toast.classList.remove("is-visible"), 2000);
}

function renderAlternatives() {
  alternativesList.replaceChildren();
  const current = captionText.value;
  const alternatives = currentCaptions.filter((caption) => caption !== current).slice(0, 3);
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

async function runGeneration(button, idleLabel) {
  setLoading(button, true, idleLabel);
  try {
    const captions = await generateCaptions(currentOptions);
    showResult(captions);
    saveCaption({
  content: captions[0],
  platform: currentOptions.platform,
  tone: currentOptions.tone
});
  } catch {
    showToast("Something went wrong. Please try again.");
  } finally {
    setLoading(button, false, idleLabel);
  }
}

form.addEventListener("submit", (event) => {
  event.preventDefault();

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

  currentOptions = readOptions();
  runGeneration(generateBtn, GENERATE_LABEL);
});

descriptionInput.addEventListener("input", () => {
  if (descriptionInput.value.trim()) {
    descriptionField.classList.remove("has-error");
    descriptionError.hidden = true;
    descriptionInput.removeAttribute("aria-invalid");
  }
});

regenerateBtn.addEventListener("click", () => {
  runGeneration(regenerateBtn, REGENERATE_LABEL);
});

editBtn.addEventListener("click", () => {
  const isEditing = captionText.readOnly;
  captionText.readOnly = !isEditing;
  editBtn.textContent = isEditing ? "Done" : "Edit";
  editBtn.setAttribute("aria-pressed", String(isEditing));
  if (isEditing) {
    captionText.focus();
    captionText.setSelectionRange(captionText.value.length, captionText.value.length);
  } else {
    renderAlternatives();
  }
});

copyBtn.addEventListener("click", async () => {
  try {
    await navigator.clipboard.writeText(captionText.value);
    showToast("Caption copied to clipboard");
  } catch {
    showToast("Couldn't copy. Select the text and copy it manually.");
  }
});
