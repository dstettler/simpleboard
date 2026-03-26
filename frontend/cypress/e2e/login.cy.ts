describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('/login');
  });

  it('loads the login page', () => {
    cy.contains('h2', 'Login').should('be.visible');
    cy.get('input[formcontrolname="email"]').should('exist');
    cy.get('input[formcontrolname="password"]').should('exist');
  });

  it('shows validation errors when fields are touched and left empty', () => {
  cy.get('input[formcontrolname="email"]').focus().blur();
  cy.get('input[formcontrolname="password"]').focus().blur();

  cy.contains('Enter a valid email').should('be.visible');
  cy.contains('Password must be at least 6 characters').should('be.visible');
});

  it('switches to register mode', () => {
    cy.get('.link-button').contains('Register').click();

    cy.contains('h2', 'Register').should('be.visible');
    cy.get('input[formcontrolname="userId"]').should('exist');
    cy.get('input[formcontrolname="name"]').should('exist');
    cy.get('input[formcontrolname="email"]').should('exist');
    cy.get('input[formcontrolname="password"]').should('exist');
  });

  it('shows success message after register', () => {
    cy.get('.link-button').contains('Register').click();

    cy.get('input[formcontrolname="userId"]').type('sreeram01');
    cy.get('input[formcontrolname="name"]').type('Sreeram');
    cy.get('input[formcontrolname="email"]').type('sreeram@test.com');
    cy.get('input[formcontrolname="password"]').type('123456');

    cy.contains('button', 'Register').click();

    cy.contains('Successfully Registered').should('be.visible');
  });
});