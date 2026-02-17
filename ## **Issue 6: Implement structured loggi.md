## **Issue 6: Implement structured logging + correlation IDs**

### **Summary**
Add structured logging to the deletion workflow, including correlation IDs for tracing Graph API calls.

### **Acceptance Criteria**
- [ ] Logs include correlation ID, path, operation type, and result.
- [ ] Logs are machine‑parseable (JSON).
- [ ] Logging respects verbosity flags.

### **Checklist**
- [ ] Add correlation ID generator
- [ ] Add structured logging wrapper
- [ ] Add logs to deletion workflow
- [ ] Add unit tests for logging behavior

### **Dependencies**
- Issue 3: Graph delete integration

### **Labels**
type: task  
component: observability  
status: ready  
priority: medium  

---

## **Issue 7: Implement error handling for missing or inaccessible files**

### **Summary**
Implement robust error handling for deletion failures, including missing files, permission issues, and transient Graph errors.

### **Acceptance Criteria**
- [ ] Missing files return a non‑zero exit code.
- [ ] Permission errors return clear messages.
- [ ] Transient errors include retry guidance.
- [ ] Errors are logged with structured metadata.

### **Checklist**
- [ ] Add error mapping layer
- [ ] Add user‑friendly error messages
- [ ] Add unit tests for error scenarios

### **Dependencies**
- Issue 3: Graph delete integration

### **Labels**
type: task  
component: cli  
status: ready  
priority: high  

---

## **Issue 8: Add unit tests for file deletion logic**

### **Summary**
Add comprehensive unit tests for the file deletion workflow.

### **Acceptance Criteria**
- [ ] Covers success, failure, and edge cases.
- [ ] Covers flag behavior (`--force`, `--dry-run`).
- [ ] Covers path resolution and error handling.

### **Checklist**
- [ ] Write unit tests for deletion executor
- [ ] Write unit tests for flag behavior
- [ ] Write unit tests for error handling

### **Dependencies**
- All previous deletion tasks

### **Labels**
type: task  
component: test  
status: ready  
priority: medium  

---

## **Issue 9: Add integration tests for file deletion**

### **Summary**
Add integration tests validating end‑to‑end deletion behavior using mock or sandbox Graph API.

### **Acceptance Criteria**
- [ ] Validates deletion of existing files.
- [ ] Validates error behavior for missing files.
- [ ] Validates flag behavior.
- [ ] Validates logging output.

### **Checklist**
- [ ] Implement integration test harness
- [ ] Add deletion integration tests
- [ ] Add logging validation

### **Dependencies**
- Issue 8: Unit tests

### **Labels**
type: task  
component: test  
status: ready  
priority: medium  

---

# 👍 Ready for the rest?

If you like this format, I will generate:

### **All remaining issues for:**
- Delete folder  
- Create folder  
- Create file  
- Upload file  
- Upload folder  
- Download file  
- Download folder  
- Read file  
- Edit file  
- Move file  
- Move folder  

Just say **“continue”** or tell me if you want them grouped, batched, or delivered as a single giant Markdown file.