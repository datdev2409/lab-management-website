# Unsaved Changes Feature - Visual Guide

## UI Elements

### 1. Status Row - Normal State (No Changes)
```
┌────────────────────────────────────────────────────────────────┐
│ Trạng thái                                                     │
├────────────────────────────────────────────────────────────────┤
│  [New]  🔓 Đang chỉnh sửa          [👁️ Xem]  [✏️ Sửa]         │
└────────────────────────────────────────────────────────────────┘
```

### 2. Status Row - With Unsaved Changes
```
┌────────────────────────────────────────────────────────────────┐
│ Trạng thái                                                     │
├────────────────────────────────────────────────────────────────┤
│  [New]  [⚠️ Chưa lưu]  🔓 Đang chỉnh sửa  [👁️ Xem] [✏️ Sửa]   │
└────────────────────────────────────────────────────────────────┘
         ↑
         │
    Yellow badge appears here when changes detected!
```

## Badge Styling
- **Background**: Yellow (`bg-warning` Bootstrap class)
- **Text Color**: Dark (default Bootstrap warning text)
- **Icon**: ⚠️ Warning triangle
- **Text**: "Chưa lưu" (Unsaved in Vietnamese)
- **Display**: `x-show="hasUnsavedChanges"`

## Complete Form Layout

```
┌──────────────────────────────────────────────────────────────────────┐
│ Tạo phiếu xét nghiệm mới                         [+ Thêm bệnh nhân] │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Trạng thái                                                     │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │  [New]  [⚠️ Chưa lưu]  🔓 Đang chỉnh sửa  [👁️ Xem] [✏️ Sửa]   │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │ Bệnh nhân                                                      │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │  [Nguyễn Văn A                        ▼]                       │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │ Bác sĩ chỉ định (tùy chọn)                                    │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │  [Bác sĩ Trần Thị B                   ▼]                       │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │ Tên gói xét nghiệm                                            │ │
│  ├────────────────────────────────────────────────────────────────┤ │
│  │  Gói xét nghiệm tổng quát        [Chọn gói xét nghiệm khác]   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ Danh sách xét nghiệm                                          │ │
│  ├────┬────────┬──────────┬────────┬────────┬────────┬───┬────┤   │
│  │Tên │Đơn giá │GT bình   │Đơn vị  │Kết quả │Kết quả │Bất│    │   │
│  │    │        │thường    │        │        │(text)  │thường│  │   │
│  ├────┼────────┼──────────┼────────┼────────┼────────┼───┼────┤   │
│  │GLU │50,000  │70-110    │mg/dL   │[95  ]  │[     ] │☐  │ ✕ │   │
│  ├────┼────────┼──────────┼────────┼────────┼────────┼───┼────┤   │
│  │HbA1c│80,000 │<5.7      │%       │[6.2 ]  │[     ] │☑⚠️│ ✕ │   │
│  │    │        │[0-5.7]   │        │        │        │   │    │   │
│  └────┴────────┴──────────┴────────┴────────┴────────┴───┴────┘   │
│                                                                      │
│  [← Trở lại]  [Lưu thay đổi (Ctrl + S)]                            │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

## Interaction Flow

### Scenario 1: User Makes Changes
```
Step 1: User opens form
Status: [New]
Badge: Hidden ❌

Step 2: User selects patient
Status: [New] [⚠️ Chưa lưu]
Badge: Visible ✅

Step 3: User clicks save
→ API call succeeds
Status: [New]
Badge: Hidden ❌
→ Success notification shown
```

### Scenario 2: User Tries to Leave Without Saving
```
Step 1: User makes changes
Status: [New] [⚠️ Chưa lưu]
Badge: Visible ✅

Step 2: User clicks browser back button
→ Browser shows native dialog:

   ┌────────────────────────────────────────────┐
   │  Leave site?                               │
   │                                            │
   │  Changes you made may not be saved.        │
   │                                            │
   │         [Cancel]         [Leave]           │
   └────────────────────────────────────────────┘

Step 3a: User clicks "Cancel"
→ Stays on page
→ Changes preserved

Step 3b: User clicks "Leave"
→ Navigates away
→ Changes lost
```

### Scenario 3: User Reloads Page
```
Step 1: User makes changes
Status: [New] [⚠️ Chưa lưu]
Badge: Visible ✅

Step 2: User presses F5 or Ctrl+R
→ Browser shows native dialog:

   ┌────────────────────────────────────────────┐
   │  Reload site?                              │
   │                                            │
   │  Changes you made may not be saved.        │
   │                                            │
   │         [Cancel]         [Reload]          │
   └────────────────────────────────────────────┘
```

## Color Reference

### Badge Colors
- **Normal State**: Gray badge (`badge text-bg-secondary`)
  - Background: #6c757d (Bootstrap secondary)
  - Text: White

- **Unsaved State**: Yellow badge (`badge text-bg-warning`)
  - Background: #ffc107 (Bootstrap warning)
  - Text: Dark
  - Icon: ⚠️

### Form States
- **View Mode**: 🔒 Green lock icon
  - Text: "Chỉ đọc" (Read-only)
  
- **Edit Mode**: 🔓 Yellow unlock icon
  - Text: "Đang chỉnh sửa" (Editing)

### Buttons
- **View Button**: Success (green)
  - Active: Solid green
  - Inactive: Outline green
  
- **Edit Button**: Warning (yellow)
  - Active: Solid yellow
  - Inactive: Outline yellow

## Responsive Behavior

### Desktop (> 992px)
```
[Badge] [Mode] [Buttons] ← All in one row
```

### Tablet (768px - 992px)
```
[Badge] [Mode]
[Buttons] ← Buttons may wrap
```

### Mobile (< 768px)
```
[Badge]
[Mode]
[Buttons] ← Stack vertically if needed
```

## Animation Notes
- **Badge Appearance**: Instant (no animation)
- **Badge Disappearance**: Instant (no animation)
- **Browser Dialog**: Native animation (varies by browser)
- **Alpine.js**: Uses `x-show` (display: none/block)

## Accessibility
- **Badge**: Visible to all users, screen readers will announce
- **Warning Icon**: Decorative, text is sufficient
- **Browser Dialog**: Fully accessible native dialog
- **Keyboard Navigation**: All interactive elements focusable

## Browser-Specific Dialog Text

### Chrome/Edge
```
Leave site?
Changes you made may not be saved.
[Cancel] [Leave]
```

### Firefox
```
This page is asking you to confirm that you want to leave
— information you've entered may not be saved.
[Leave Page] [Stay on Page]
```

### Safari
```
Are you sure you want to leave this page?
[Cancel] [Leave]
```

Note: Exact text may vary by browser version and language settings.

## Implementation Notes

### Badge Visibility Logic
```javascript
x-show="hasUnsavedChanges"
```
- Uses Alpine.js `x-show` directive
- Reactive: Updates immediately when `hasUnsavedChanges` changes
- No manual DOM manipulation required

### Badge HTML Structure
```html
<span x-show="hasUnsavedChanges" class="badge text-bg-warning">
    <i class="bi bi-exclamation-triangle-fill"></i>
    Chưa lưu
</span>
```
- Bootstrap 5 badge component
- Bootstrap Icons for warning icon
- Vietnamese text for local users

### Integration with Existing UI
- Fits seamlessly in existing status row
- Uses same Bootstrap styling as other badges
- Maintains visual hierarchy
- No layout shifts (space is preserved)

## Testing Checklist

Visual Testing:
- [ ] Badge appears immediately on field change
- [ ] Badge has correct yellow color
- [ ] Badge has warning icon
- [ ] Badge text is readable
- [ ] Badge disappears after save
- [ ] Layout doesn't shift when badge appears/disappears

Interaction Testing:
- [ ] Browser warning shows on back button
- [ ] Browser warning shows on reload
- [ ] Browser warning shows on tab close
- [ ] No warning when no changes made
- [ ] Save button clears warning

Cross-browser Testing:
- [ ] Chrome: Badge + warning work
- [ ] Firefox: Badge + warning work
- [ ] Safari: Badge + warning work
- [ ] Edge: Badge + warning work

Responsive Testing:
- [ ] Desktop: Badge displays correctly
- [ ] Tablet: Badge displays correctly
- [ ] Mobile: Badge displays correctly
