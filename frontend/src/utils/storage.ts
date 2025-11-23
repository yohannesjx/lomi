import { Platform } from 'react-native';

// Only import SecureStore on native platforms
let SecureStore: typeof import('expo-secure-store') | null = null;
if (Platform.OS !== 'web') {
    try {
        SecureStore = require('expo-secure-store');
    } catch (e) {
        // SecureStore not available
    }
}

/**
 * Cross-platform storage utility
 * Uses SecureStore on native platforms and localStorage on web
 */
class Storage {
    async getItem(key: string): Promise<string | null> {
        if (Platform.OS === 'web') {
            try {
                return localStorage.getItem(key);
            } catch (error) {
                console.error('localStorage.getItem error:', error);
                return null;
            }
        } else {
            if (!SecureStore) {
                console.warn('SecureStore not available');
                return null;
            }
            try {
                return await SecureStore.getItemAsync(key);
            } catch (error) {
                console.error('SecureStore.getItemAsync error:', error);
                return null;
            }
        }
    }

    async setItem(key: string, value: string): Promise<void> {
        if (Platform.OS === 'web') {
            try {
                localStorage.setItem(key, value);
            } catch (error) {
                console.error('localStorage.setItem error:', error);
                throw error;
            }
        } else {
            if (!SecureStore) {
                throw new Error('SecureStore not available');
            }
            try {
                await SecureStore.setItemAsync(key, value);
            } catch (error) {
                console.error('SecureStore.setItemAsync error:', error);
                throw error;
            }
        }
    }

    async removeItem(key: string): Promise<void> {
        if (Platform.OS === 'web') {
            try {
                localStorage.removeItem(key);
            } catch (error) {
                console.error('localStorage.removeItem error:', error);
            }
        } else {
            if (!SecureStore) {
                return;
            }
            try {
                await SecureStore.deleteItemAsync(key);
            } catch (error) {
                console.error('SecureStore.deleteItemAsync error:', error);
            }
        }
    }
}

export const storage = new Storage();

