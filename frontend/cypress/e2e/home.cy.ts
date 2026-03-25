describe('Home Page', () => {
  it('loads the home page', () => {
    cy.visit('/');
    cy.contains('SimpleBoard');
  });

  it('has navigation links', () => {
    cy.visit('/');
    cy.contains('Home');
    cy.contains('Login');
  });

  it('goes to login page from navbar', () => {
    cy.visit('/');
    cy.contains('Login').click();
    cy.url().should('include', '/login');
  });
});