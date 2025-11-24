module.exports = function (api) {
    api.cache(false); // Disable cache for development
    return {
        presets: ['babel-preset-expo'],
        plugins: [
            'react-native-reanimated/plugin',
        ],
    };
};
