import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, TextInput } from 'react-native';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { ETHIOPIAN_CITIES, EthiopianCity } from '../../constants/ethiopianData';

interface CityPickerProps {
    selectedCity: EthiopianCity | null;
    onSelect: (city: EthiopianCity) => void;
}

export const CityPicker: React.FC<CityPickerProps> = ({ selectedCity, onSelect }) => {
    const [searchQuery, setSearchQuery] = useState('');
    const [isExpanded, setIsExpanded] = useState(false);

    const filteredCities = ETHIOPIAN_CITIES.filter(city =>
        city.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        city.region.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const handleSelect = (city: EthiopianCity) => {
        onSelect(city);
        setIsExpanded(false);
        setSearchQuery('');
    };

    return (
        <View style={styles.container}>
            <Text style={styles.label}>City</Text>

            <TouchableOpacity
                style={styles.selector}
                onPress={() => setIsExpanded(!isExpanded)}
            >
                <Text style={selectedCity ? styles.selectedText : styles.placeholderText}>
                    {selectedCity ? selectedCity.name : 'Select your city'}
                </Text>
                <Text style={styles.arrow}>{isExpanded ? '▲' : '▼'}</Text>
            </TouchableOpacity>

            {isExpanded && (
                <View style={styles.dropdown}>
                    <TextInput
                        style={styles.searchInput}
                        placeholder="Search cities..."
                        placeholderTextColor={COLORS.textTertiary}
                        value={searchQuery}
                        onChangeText={setSearchQuery}
                        autoFocus
                    />

                    <ScrollView style={styles.cityList} nestedScrollEnabled>
                        {filteredCities.map((city, index) => (
                            <TouchableOpacity
                                key={index}
                                style={styles.cityItem}
                                onPress={() => handleSelect(city)}
                            >
                                <Text style={styles.cityName}>{city.name}</Text>
                                <Text style={styles.cityRegion}>{city.region}</Text>
                            </TouchableOpacity>
                        ))}
                        {filteredCities.length === 0 && (
                            <Text style={styles.noResults}>No cities found</Text>
                        )}
                    </ScrollView>
                </View>
            )}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        marginBottom: SPACING.l,
    },
    label: {
        fontSize: 14,
        fontWeight: '500',
        color: COLORS.textSecondary,
        marginBottom: SPACING.s,
    },
    selector: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    selectedText: {
        fontSize: 16,
        color: COLORS.textPrimary,
    },
    placeholderText: {
        fontSize: 16,
        color: COLORS.textTertiary,
    },
    arrow: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    dropdown: {
        marginTop: SPACING.s,
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
        maxHeight: 300,
        overflow: 'hidden',
    },
    searchInput: {
        padding: SPACING.m,
        fontSize: 16,
        color: COLORS.textPrimary,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
    },
    cityList: {
        maxHeight: 250,
    },
    cityItem: {
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
    },
    cityName: {
        fontSize: 16,
        color: COLORS.textPrimary,
        marginBottom: 2,
    },
    cityRegion: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    noResults: {
        padding: SPACING.l,
        textAlign: 'center',
        color: COLORS.textSecondary,
    },
});
