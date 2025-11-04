# Unsaved Changes Detection Feature - Testing Guide

## Feature Overview
This feature detects unsaved changes in record create/edit pages and warns users before they navigate away or reload the page.

## Implementation Details

### Files Modified
1. `internal/templates/partials/record_create_form.templ`
   - Added state tracking variables
   - Added change detection logic
   - Added beforeunload event listener
   - Added unsaved badge display

2. `internal/templates/partials/record_test_table.templ`
   - Integrated change detection into input fields

### How It Works

#### State Management
- **initialFormState**: Stores JSON string of the form's initial state
- **hasUnsavedChanges**: Boolean flag tracking if changes exist
- **captureInitialState()**: Captures current form state as baseline
- **getCurrentFormState()**: Returns normalized current state for comparison
- **checkForChanges()**: Compares current state with initial state

#### Change Detection Triggers
Changes are detected on:
1. Patient selection/change
2. Doctor selection/change
3. Combo selection/change
4. Test addition/removal
5. Test result input (numeric)
6. Test result text input
7. Abnormal checkbox toggle

#### User Notifications
1. **Visual Badge**: Yellow "Chưa lưu" (Unsaved) badge appears in status row
2. **Browser Warning**: Native browser confirmation dialog when user tries to:
   - Navigate to another page
   - Reload the page
   - Close the tab/window

#### State Reset
Unsaved changes flag is cleared after:
1. Successful form submission (create or update)
2. Initial state recapture after save

## Manual Testing Instructions

### Test Case 1: Create New Record with Unsaved Changes
1. Navigate to `/phieu-xet-nghiem/new` (Create New Record page)
2. Select a patient from the dropdown
3. Verify "Chưa lưu" badge appears in the status row
4. Try to navigate away (click browser back button or enter different URL)
5. Verify browser shows "Leave site?" confirmation dialog
6. Click "Cancel" to stay on the page
7. Click "Lưu" (Save) button
8. Verify badge disappears after successful save
9. Try to navigate away - should work without warning

### Test Case 2: Edit Existing Record with Changes
1. Navigate to an existing record detail page `/phieu-xet-nghiem/{id}`
2. Click "Sửa" (Edit) button to enter edit mode
3. Change a test result value
4. Verify "Chưa lưu" badge appears
5. Try to reload the page (Ctrl+R or F5)
6. Verify browser shows confirmation dialog
7. Click "Cancel" to stay
8. Click "Lưu thay đổi" (Save Changes)
9. Verify badge disappears
10. Reload page - should work without warning

### Test Case 3: Multiple Field Changes
1. Go to create/edit record page
2. Make multiple changes:
   - Change patient
   - Select a combo
   - Edit test results
   - Toggle abnormal checkboxes
3. Verify badge shows after each change
4. Save the form
5. Verify badge clears and no warning on navigation

### Test Case 4: No Changes Scenario
1. Open existing record in edit mode
2. Do NOT make any changes
3. Try to navigate away
4. Should NOT show warning (no unsaved changes)

### Test Case 5: Keyboard Shortcuts
1. Open record form
2. Make changes (badge appears)
3. Press Ctrl+S (or Cmd+S on Mac) to save
4. Verify form saves and badge disappears

## Expected Behavior

### Success Criteria
✅ Badge appears immediately when any field is changed
✅ Browser warning shows on navigation/reload when changes exist
✅ Badge disappears after successful save
✅ No warning appears when no changes are made
✅ Works for both create and edit flows
✅ State persists across edit mode toggles

### Edge Cases Handled
- New records start with empty state
- Existing records load initial state from server
- Manual abnormal override flag changes are tracked
- Test additions and removals are detected
- Doctor/Patient/Combo changes are all tracked

## Browser Compatibility
The `beforeunload` event is supported in:
- Chrome/Edge: Shows generic message
- Firefox: Shows generic message
- Safari: Shows generic message

Note: Modern browsers ignore custom messages for security reasons and show their own standard text.

## Debugging

### Console Logs
The feature includes debug logging:
- "Initial form state captured" - when baseline state is set
- "Unsaved changes detected" - when changes are found

### Check State
In browser console, access the component state:
```javascript
Alpine.$data(document.querySelector('[x-data*="recordForm"]'))
```

This returns the component object where you can inspect:
- `hasUnsavedChanges`
- `initialFormState`
- Current form fields
