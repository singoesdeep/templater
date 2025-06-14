# Templater Troubleshooting Guide

## Common Issues and Solutions

### Template Processing Issues

#### 1. Template Not Found
**Symptoms:**
- Error: "template file not found"
- Error: "no such file or directory"

**Solutions:**
1. Verify the template path is correct
2. Check file permissions
3. Ensure the template file exists
4. Use absolute paths if needed

#### 2. Invalid Template Syntax
**Symptoms:**
- Error: "template syntax error"
- Error: "unexpected token"

**Solutions:**
1. Check for missing closing tags
2. Verify template variable names
3. Ensure proper Go template syntax
4. Use template validation tools

#### 3. Missing Variables
**Symptoms:**
- Error: "undefined variable"
- Empty output for variables

**Solutions:**
1. Check data file for required variables
2. Verify variable names match template
3. Provide default values
4. Use data validation

### Data Processing Issues

#### 1. Data File Not Found
**Symptoms:**
- Error: "data file not found"
- Error: "cannot read data file"

**Solutions:**
1. Verify data file path
2. Check file permissions
3. Ensure file exists
4. Use absolute paths if needed

#### 2. Invalid Data Format
**Symptoms:**
- Error: "invalid JSON/YAML"
- Error: "cannot parse data"

**Solutions:**
1. Validate JSON/YAML syntax
2. Check data structure
3. Verify data types
4. Use data validation tools

#### 3. Missing Data Fields
**Symptoms:**
- Error: "missing required field"
- Empty output for fields

**Solutions:**
1. Check data file for required fields
2. Verify field names
3. Provide default values
4. Use data validation

### File Operation Issues

#### 1. Permission Denied
**Symptoms:**
- Error: "permission denied"
- Error: "cannot write file"

**Solutions:**
1. Check file permissions
2. Verify directory permissions
3. Run with appropriate privileges
4. Use allowed directories

#### 2. Disk Space Issues
**Symptoms:**
- Error: "no space left on device"
- Error: "cannot write file"

**Solutions:**
1. Free up disk space
2. Check available space
3. Clean up temporary files
4. Monitor disk usage

#### 3. File Lock Issues
**Symptoms:**
- Error: "file is locked"
- Error: "cannot access file"

**Solutions:**
1. Close other programs using the file
2. Wait for file to be released
3. Use file locking mechanisms
4. Implement retry logic

### Performance Issues

#### 1. Slow Processing
**Symptoms:**
- Long processing times
- High CPU usage

**Solutions:**
1. Optimize template size
2. Use concurrent processing
3. Enable caching
4. Monitor resource usage

#### 2. High Memory Usage
**Symptoms:**
- Out of memory errors
- Slow performance

**Solutions:**
1. Monitor memory usage
2. Clean up resources
3. Optimize data structures
4. Use garbage collection

#### 3. File I/O Bottlenecks
**Symptoms:**
- Slow file operations
- High disk usage

**Solutions:**
1. Use buffered I/O
2. Implement caching
3. Optimize file operations
4. Monitor disk usage

### Security Issues

#### 1. Template Security
**Symptoms:**
- Error: "dangerous operation detected"
- Error: "invalid template content"

**Solutions:**
1. Validate template content
2. Use secure templates
3. Implement security checks
4. Follow security guidelines

#### 2. Data Security
**Symptoms:**
- Error: "invalid data format"
- Error: "security violation"

**Solutions:**
1. Validate input data
2. Sanitize user input
3. Use secure data handling
4. Follow security guidelines

#### 3. File Security
**Symptoms:**
- Error: "invalid file path"
- Error: "security violation"

**Solutions:**
1. Validate file paths
2. Use secure file operations
3. Implement access control
4. Follow security guidelines

### Plugin Issues

#### 1. Plugin Loading
**Symptoms:**
- Error: "cannot load plugin"
- Error: "plugin not found"

**Solutions:**
1. Verify plugin path
2. Check plugin compatibility
3. Ensure plugin is compiled
4. Use correct plugin format

#### 2. Plugin Execution
**Symptoms:**
- Error: "plugin execution failed"
- Error: "plugin error"

**Solutions:**
1. Check plugin code
2. Verify plugin interface
3. Test plugin separately
4. Check plugin logs

#### 3. Plugin Compatibility
**Symptoms:**
- Error: "incompatible plugin"
- Error: "version mismatch"

**Solutions:**
1. Check plugin version
2. Update plugin
3. Verify compatibility
4. Use compatible versions

## Debugging Tips

### 1. Enable Debug Logging
```bash
templater --debug generate -t template.tmpl -d data.json
```

### 2. Check File Permissions
```bash
ls -l template.tmpl data.json
```

### 3. Validate Template
```bash
templater validate template.tmpl
```

### 4. Check Data Format
```bash
templater validate-data data.json
```

### 5. Monitor Resources
```bash
templater --monitor generate -t template.tmpl -d data.json
```

## Getting Help

### 1. Check Documentation
- Read the user guide
- Review API documentation
- Check best practices
- Look for examples

### 2. Search Issues
- Check GitHub issues
- Look for similar problems
- Review solutions
- Check closed issues

### 3. Community Support
- Ask on forums
- Join discussions
- Share experiences
- Get help from community

### 4. Report Issues
- Use issue templates
- Provide details
- Include logs
- Share steps to reproduce 