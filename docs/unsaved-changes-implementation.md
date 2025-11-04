# Unsaved Changes Detection Feature - Implementation Summary

## Overview
This feature adds comprehensive unsaved changes detection to laboratory test record create/edit pages, preventing accidental data loss when users navigate away without saving.

## Problem Statement
Users could lose work if they:
- Navigate to another page without saving
- Reload the page
- Close the browser tab
- Close the browser window

## Solution
Implemented a robust change detection system with dual-layer user notifications:
1. **Visual feedback**: Warning badge in the UI
2. **Browser protection**: Native confirmation dialog before leaving

## Technical Implementation

### Architecture
The solution uses Alpine.js reactive state management with JSON-based state comparison:

```javascript
// State properties
initialFormState: null,      // JSON string of baseline state
hasUnsavedChanges: false,    // Boolean flag for UI/warnings

// Core methods
captureInitialState()        // Saves current state as baseline
getCurrentFormState()        // Returns normalized state object
checkForChanges()           // Compares current vs baseline
```

### State Tracked
All form fields are monitored:
- **Patient**: ID and name
- **Doctor**: ID and name (optional)
- **Combo**: Selected combo name
- **Test Results**: For each test:
  - ID and name
  - Numeric result value
  - Text result value
  - Abnormal flag
  - Manual abnormal override flag

### Change Detection Flow
```
User Action → Field Update → checkForChanges() → 
Compare States → Update hasUnsavedChanges → 
Show/Hide Badge + Enable/Disable Browser Warning
```

### Event Listeners
1. **Form field changes**: All inputs trigger `checkForChanges()`
2. **beforeunload**: Browser event prevents navigation when changes exist
3. **Submit success**: Clears unsaved flag and recaptures baseline

## User Experience

### Visual Indicators
- **Badge Location**: Status row (first row in form table)
- **Badge Style**: Yellow background with warning icon
- **Badge Text**: "Chưa lưu" (Vietnamese for "Unsaved")
- **Badge Behavior**: 
  - Shows immediately when any field changes
  - Hides immediately after successful save

### Browser Warnings
- **Trigger**: User attempts to navigate away or reload
- **Condition**: Only when `hasUnsavedChanges === true`
- **Dialog Type**: Native browser confirmation (cannot be customized)
- **Options**: User can choose to stay or leave

### Keyboard Shortcuts
- **Ctrl+S / Cmd+S**: Save form (existing feature)
- **Ctrl+E / Cmd+E**: Toggle edit mode (existing feature)

## Code Quality

### Testing
- ✅ Compiles without errors
- ✅ Templates generate correctly
- ✅ Passes go fmt
- ✅ Passes go vet
- ✅ No CodeQL security alerts
- ✅ Code review feedback addressed

### Documentation
- Comprehensive testing guide in `docs/unsaved-changes-feature-test.md`
- Inline code comments explaining design decisions
- Clear variable naming for maintainability

### Performance
- **State comparison**: O(n) where n is number of tests
- **Memory usage**: One JSON string per form instance
- **Event overhead**: Minimal, only one beforeunload listener per page
- **UI updates**: Reactive, no manual DOM manipulation

## Browser Compatibility
Works on all modern browsers:
- Chrome/Edge ✓
- Firefox ✓
- Safari ✓

Note: Browser confirmation message text cannot be customized for security reasons (all browsers show their standard warning).

## Edge Cases Handled

### ✅ New Record Creation
- Initial state captured after Alpine.js initialization
- Empty fields don't trigger false positives

### ✅ Existing Record Loading
- Initial state captured after data fetch
- Loading state doesn't trigger warnings

### ✅ Edit Mode Toggling
- View mode doesn't trigger change detection
- Edit mode properly tracks all changes

### ✅ Save Success
- State recaptured after successful API call
- Badge clears immediately
- Navigation allowed without warning

### ✅ Save Failure
- Unsaved flag remains true
- User can retry save
- Navigation still protected

### ✅ Multiple Rapid Changes
- Debounced through Alpine.js reactivity
- Only final state matters

### ✅ Form Field Edge Cases
- Null/undefined patient or doctor
- Empty test results array
- Missing optional fields
- Non-numeric test result values

## Limitations

### By Design
1. **No undo/redo**: Feature only detects changes, doesn't track history
2. **No autosave**: Intentionally requires explicit user action
3. **No custom warning text**: Browser security restriction

### Technical
1. **Requires JavaScript**: Won't work if JavaScript is disabled
2. **Single page only**: State doesn't persist across page reloads
3. **Client-side only**: No server-side session management

## Maintenance Notes

### Future Enhancements
Consider adding:
- Autosave draft functionality
- Change history/undo capability
- More granular change notifications (which field changed)
- Analytics on save frequency and abandoned edits

### Common Issues
If the feature stops working, check:
1. Alpine.js is loaded and initialized
2. `recordForm` component is properly instantiated
3. Browser console for JavaScript errors
4. Template generation is up to date

### Debugging
Access component state in browser console:
```javascript
Alpine.$data(document.querySelector('[x-data*="recordForm"]'))
```

Check for:
- `hasUnsavedChanges` value
- `initialFormState` content
- Current form field values

## Files Changed
1. `internal/templates/partials/record_create_form.templ`
   - Added state management properties
   - Added change detection methods
   - Added beforeunload event listener
   - Added unsaved badge to UI
   - Integrated checkForChanges() calls

2. `internal/templates/partials/record_test_table.templ`
   - Added checkForChanges() to input event handlers
   - Added checkForChanges() to checkbox change handlers

3. `docs/unsaved-changes-feature-test.md`
   - Comprehensive manual testing guide
   - Browser compatibility notes
   - Debugging instructions

## Security Considerations
- ✅ No XSS vulnerabilities (all data properly handled by Alpine.js)
- ✅ No sensitive data in localStorage or sessionStorage
- ✅ No external API calls for change detection
- ✅ CodeQL scan passed with 0 alerts
- ✅ State comparison doesn't expose sensitive info

## Performance Impact
- **Minimal**: Only JSON serialization on field changes
- **No API calls**: All comparison done client-side
- **No polling**: Event-driven only
- **Small memory footprint**: Single JSON string per form

## Accessibility
- Visual badge is visible to sighted users
- Browser warning is accessible to screen readers
- No keyboard trap or focus issues introduced
- Form remains fully keyboard-navigable

## Conclusion
This implementation provides robust protection against accidental data loss while maintaining excellent performance and user experience. The solution is maintainable, well-documented, and follows Alpine.js best practices.
