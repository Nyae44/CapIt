const STORAGE_KEY = "capit_saved_captions";

// Get all saved captions
export function getCaptions() {
    return JSON.parse(
        localStorage.getItem(STORAGE_KEY) || "[]"
    );
}

// Save a caption
export function saveCaption(caption) {
    const captions = getCaptions();

    const newCaption = {
        id: crypto.randomUUID(),
        ...caption,
        createdAt: new Date().toISOString()
    };

    captions.unshift(newCaption);

    localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify(captions)
    );

    return newCaption;
}

// Delete a caption
export function deleteCaption(id) {
    const captions = getCaptions();

    const updatedCaptions = captions.filter(
        caption => caption.id !== id
    );

    localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify(updatedCaptions)
    );
}

// Clear all captions
export function clearCaptions() {
    localStorage.removeItem(STORAGE_KEY);
}

