# Planning & Design for Atala PRISM Integration
1. Research Atala PRISM API **(Complete)**
   - Review API documentation of endpoints and request/response types
   - Particular attention to Verified Credentials & Proofs
   - Launch prism agents locally in docker
   - Develop simple client to audition API endpoints
2. Research ProofSpace **(Complete)**
   - Wallet usage and QR Code UI integration
   - Interaction features & parameters, webhooks, and UX flow
3. Identify sundae-governance integration points **(Complete)**

## Development of Atala PRISM Integration
1. Modify sundae-governance types to allow proposal configuration, permissions, and weights based on verified credential schemas **(Complete)**
    - Voting Permissions:
      - VC Schema URL / ProofSpace Interaction ID
      - VC field
      - Formula (e.g. `count`, `sum`)
      - Value required to allow voting
    - Voting Weights: same as above with value representing vote weight
2. Modify Sundae Governance UI to allow customers to select new proposal type (`PROOFSPACE_INTERACTION_POC`) and enter permission and weight values **(Complete)**
3. Implement voter registration flow **(Complete)**
   - Create interaction in ProofSpace backend which filter VC's by issuer DID and VC Schema
   - Integrate ProofSpace wallet with UI to allow customer's users to choose a VC(s) which satisfies proposal permission requirements
   - Create webhook to receive and validate VC from ProofSpace
   - Verify VC satisfies proposal permission requirements
   - Vote permission is stored on-chain
4. Implement vote casting flow **(Complete)**
   - Registered users receive push notification when voting opens
   - Customer's users to choose a VC(s) which satisfies proposal weight requirements
   - Vote weight is calculated for storage on-chain

## Testing, Debugging and User Testing
1. Finalize acceptance criteria for customer and user flows **(In-Progress)**
2. Create mocks and stubs for verified credential interactions to test webhooks **(Complete)**
3. Bring unit and integration test coverage of acceptance criteria to 100% **(In-Progress)**

## Pilot Launch of Atala PRISM Integration
1. ProofSpace integration removes dependency on PRISM infrastructure
2. Production release per sundae-governance product roadmap

