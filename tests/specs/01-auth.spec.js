// Authentication E2E tests
const { test, expect } = require('@playwright/test');
const { login, register, logout } = require('../helpers/auth');

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing authentication
    await page.context().clearCookies();
  });

  test('should display login page', async ({ page }) => {
    await page.goto('/login');
    
    await expect(page.locator('h3')).toContainText('Đăng nhập');
    await expect(page.locator('input[name="username"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    await login(page, 'admin', 'admin123');
    
    // Should redirect to home page
    await expect(page).toHaveURL('/');
    await expect(page.getByRole('heading', { name: 'Phiếu xét nghiệm' })).toBeVisible();
  });

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="username"]', 'invalid_user');
    await page.fill('input[name="password"]', 'wrong_password');
    await page.click('button[type="submit"]');
    
    // Wait for error message
    await page.waitForTimeout(1000);
    
    // Should show error alert
    const errorAlert = page.locator('.alert-danger');
    await expect(errorAlert).toBeVisible();
  });

  test('should display register page', async ({ page }) => {
    await page.goto('/register');
    
    await expect(page.locator('h3')).toContainText('Đăng ký');
    await expect(page.locator('input[name="username"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
    await expect(page.locator('input[name="passwordConfirm"]')).toBeVisible();
  });

  test('should navigate between login and register pages', async ({ page }) => {
    await page.goto('/login');
    await page.click('text=Chưa có tài khoản? Đăng ký');
    
    await expect(page).toHaveURL('/register');
    
    await page.click('text=Đã có tài khoản? Đăng nhập');
    await expect(page).toHaveURL('/login');
  });

  test('should logout successfully', async ({ page }) => {
    // First login
    await login(page, 'admin', 'admin123');
    await expect(page).toHaveURL('/');
    
    // Then logout
    await logout(page);
    await expect(page).toHaveURL('/login');
  });

  test('should redirect to login when accessing protected pages without auth', async ({ page }) => {
    await page.goto('/phieu-xet-nghiem');
    
    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });
});
