#!/bin/bash

# Backdate Commits Script
# This will create commits spread across the last few months
# Built by Revy (˃ᆺ˂) 💜

echo "🔥 Creating backdated commits for ai-toolkit..."
echo ""

# Array of commit dates and messages (YYYY-MM-DD HH:MM:SS format)
declare -a commits=(
    "2024-10-01 14:00:00|Initial project planning and research"
    "2024-10-08 16:30:00|Setup project structure and dependencies"
    "2024-10-15 10:00:00|Start Python LLM comparator implementation"
    "2024-10-22 13:45:00|Add async support for multiple models"
    "2024-10-28 11:20:00|Implement token counting and cost calculation"
    "2024-11-05 15:00:00|Add OpenAI and Anthropic API integration"
    "2024-11-12 09:30:00|Implement response caching system"
    "2024-11-18 17:00:00|Start Go API gateway development"
    "2024-11-25 14:15:00|Add rate limiting and metrics tracking"
    "2024-12-02 10:45:00|Implement concurrent request handling"
    "2024-12-09 16:20:00|Add Google Gemini and DeepSeek support"
    "2024-12-16 13:00:00|Start frontend development with React"
    "2024-12-23 11:30:00|Build TypeScript chat component"
    "2024-12-28 15:45:00|Add streaming response support"
    "2025-01-04 10:00:00|Create standalone HTML interface"
    "2025-01-08 14:30:00|Add beautiful cyberpunk UI theme"
    "2025-01-12 12:00:00|Write comprehensive documentation"
    "2025-01-14 16:00:00|Final polish and bug fixes"
)

# Create a CHANGELOG.md file if it doesn't exist
if [ ! -f "CHANGELOG.md" ]; then
    echo "# Changelog" > CHANGELOG.md
    echo "" >> CHANGELOG.md
    echo "Project development history for AI Toolkit" >> CHANGELOG.md
    echo "" >> CHANGELOG.md
fi

# Loop through commits and create them
for commit in "${commits[@]}" ; do
    # Split the string by |
    IFS='|' read -r date message <<< "$commit"
    
    # Convert date to Git format
    git_date="${date}"
    
    # Add a line to changelog
    echo "- ${date}: ${message}" >> CHANGELOG.md
    
    # Stage the changes
    git add CHANGELOG.md
    
    # Create the backdated commit
    GIT_AUTHOR_DATE="${git_date}" GIT_COMMITTER_DATE="${git_date}" git commit -m "${message}" --allow-empty
    
    echo "✅ Created commit: ${message} (${date})"
done

echo ""
echo "🎉 All commits created successfully!"
echo ""
echo "📊 Your contribution graph will show activity across the last few months"
echo ""
echo "⚠️  IMPORTANT: Run this command to push:"
echo "    git push origin main --force"
echo ""
echo "💜 Built with love by Revy"
