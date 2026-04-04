# Configuration Management

The OneDrive CLI (odc) needs a flexible and predictable way to understand how it should behave. Users may run odc in many different environments—local machines, CI pipelines, containers, or ephemeral shells—and each of those environments may require different configuration sources. To support this, odc uses a configuration system designed around three core principles:

1. **Configurations can come from multiple places**  
2. **Configurations should be easy to override**  
3. **Configurations should not assume a specific file format**

This section explains the reasoning behind the configuration system and how it fits into the broader architecture of odc.

## Multiple Sources of Configuration

odc supports configuration from several locations, each serving a different purpose:

### **1. CLI‑Provided Configuration**
Users can provide configuration directly on the command line:

```shell
odc --config /path/to/config.json
```

This is ideal for:

- temporary overrides  
- CI/CD pipelines  
- testing new settings  
- running odc without modifying local files  

CLI‑provided configuration always takes precedence because it represents the user’s most explicit intent.

### **2. Installed Configuration**
When no CLI override is provided, odc falls back to the user’s installed configuration. This is stored in OS‑appropriate locations:

- Linux: `~/.config/odc/`
- macOS: `~/Library/Application Support/odc/`
- Windows: `%APPDATA%\odc\`

Installed configuration represents the user’s long‑term preferences and is the foundation for profile management.

## File‑Type Agnostic by Design

odc does not assume that configuration files are JSON. While JSON is the default today, the configuration manager is intentionally designed to be **file‑type agnostic**.

This means:

- JSON, YAML, TOML, or other formats can be supported  
- encrypted configuration files can be added later  
- enterprise or remote configuration sources can be plugged in  
- users are not locked into a single format  

This flexibility is achieved through a **Loader abstraction**, which allows odc to read configuration files without caring about their underlying format.

## Lazy Loading and Caching

Configurations are loaded **only when needed**. This keeps odc fast and avoids unnecessary disk access.

The process works like this:

1. odc asks the configuration manager for a configuration by name  
2. the configuration manager checks its cache  
3. if not cached, it loads the configuration from disk using the appropriate loader  
4. the loaded configuration is cached for future use  

This approach ensures:

- fast repeated access  
- minimal I/O  
- predictable behavior even with many profiles  

## Separation of Concerns

The configuration system is intentionally separated from authentication and profile management.

- **ConfigurationService** loads and caches configuration files  
- **CacheService** stores authentication records  
- **ProfileService** combines configuration + identity + overrides  

This separation keeps each component simple and focused, and it allows odc to evolve without breaking existing behavior.

## Why This Matters

A flexible configuration system is essential for a CLI tool that needs to work across:

- personal machines
- automated systems
- multiple profiles
- multiple authentication methods

By keeping the configuration manager format‑agnostic, lazy‑loaded, and decoupled from other services, odc ensures that users can adapt it to their needs without friction.
