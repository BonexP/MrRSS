/**
 * Centralized default values for settings
 * These values are loaded from the shared config/defaults.json file
 * to ensure consistency between frontend and backend.
 */

// Import the shared defaults from the config directory
import sharedDefaults from '../../../config/defaults.json';

export const settingsDefaults = sharedDefaults;

// Type for the defaults object
export type SettingsDefaults = typeof settingsDefaults;
