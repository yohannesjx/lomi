import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '../../components/ui/Button';
import { COLORS, SPACING, SIZES } from '../../theme/colors';

const INTERESTS = [
    { id: 'buna', label: '☕ Buna Lover', category: 'Culture' },
    { id: 'music', label: '🎵 Music', category: 'Art' },
    { id: 'travel', label: '✈️ Travel', category: 'Lifestyle' },
    { id: 'movies', label: '🎬 Movies', category: 'Art' },
    { id: 'fitness', label: '💪 Fitness', category: 'Lifestyle' },
    { id: 'foodie', label: '🍲 Foodie', category: 'Lifestyle' },
    { id: 'tech', label: '💻 Tech', category: 'Work' },
    { id: 'art', label: '🎨 Art', category: 'Art' },
    { id: 'faith', label: '🙏 Faith', category: 'Values' },
    { id: 'reading', label: '📚 Reading', category: 'Hobbies' },
    { id: 'dancing', label: '💃 Dancing', category: 'Hobbies' },
    { id: 'football', label: '⚽ Football', category: 'Sports' },
    { id: 'photography', label: '📸 Photography', category: 'Art' },
    { id: 'fashion', label: '👗 Fashion', category: 'Lifestyle' },
    { id: 'nature', label: '🌿 Nature', category: 'Lifestyle' },
];

export const InterestsScreen = ({ navigation }: any) => {
    const [selectedInterests, setSelectedInterests] = useState<string[]>([]);

    const toggleInterest = (id: string) => {
        if (selectedInterests.includes(id)) {
            setSelectedInterests(selectedInterests.filter(i => i !== id));
        } else {
            if (selectedInterests.length >= 5) {
                alert('You can only select up to 5 interests');
                return;
            }
            setSelectedInterests([...selectedInterests, id]);
        }
    };

    const handleFinish = () => {
        // Navigate to Gender Preference screen
        navigation.navigate('GenderPreference');
    };

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.content}>
                <View style={styles.header}>
                    <Text style={styles.stepIndicator}>Step 3 of 4</Text>
                    <Text style={styles.title}>What are you into?</Text>
                    <Text style={styles.subtitle}>Pick up to 5 interests to find better matches</Text>
                </View>

                <ScrollView contentContainerStyle={styles.tagsContainer}>
                    {INTERESTS.map((interest) => {
                        const isSelected = selectedInterests.includes(interest.id);
                        return (
                            <TouchableOpacity
                                key={interest.id}
                                style={[
                                    styles.tag,
                                    isSelected && styles.tagSelected
                                ]}
                                onPress={() => toggleInterest(interest.id)}
                            >
                                <Text style={[
                                    styles.tagLabel,
                                    isSelected && styles.tagLabelSelected
                                ]}>
                                    {interest.label}
                                </Text>
                            </TouchableOpacity>
                        );
                    })}
                </ScrollView>

                <View style={styles.footer}>
                    <Text style={styles.counter}>
                        {selectedInterests.length}/5 selected
                    </Text>
                    <Button
                        title="Finish Profile"
                        onPress={handleFinish}
                        disabled={selectedInterests.length < 3}
                        size="large"
                    />
                </View>
            </View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    content: {
        flex: 1,
        padding: SPACING.l,
    },
    header: {
        marginBottom: SPACING.xl,
    },
    stepIndicator: {
        color: COLORS.primary,
        fontWeight: 'bold',
        marginBottom: SPACING.s,
    },
    title: {
        fontSize: 28,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    subtitle: {
        fontSize: 16,
        color: COLORS.textSecondary,
    },
    tagsContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        gap: SPACING.s,
        paddingBottom: SPACING.xl,
    },
    tag: {
        paddingHorizontal: SPACING.m,
        paddingVertical: SPACING.s,
        borderRadius: SIZES.radiusXL,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
    },
    tagSelected: {
        borderColor: COLORS.primary,
        backgroundColor: COLORS.primary,
    },
    tagLabel: {
        color: COLORS.textSecondary,
        fontSize: 16,
        fontWeight: '500',
    },
    tagLabelSelected: {
        color: COLORS.background,
        fontWeight: 'bold',
    },
    footer: {
        marginTop: 'auto',
        paddingTop: SPACING.m,
    },
    counter: {
        color: COLORS.textSecondary,
        textAlign: 'center',
        marginBottom: SPACING.m,
    },
});
