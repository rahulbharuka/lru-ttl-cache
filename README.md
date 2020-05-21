# LRU TTL Cache Library

This project implements an in-memory Least Recently Used (LRU) and Time To Live (TTL) based cache library to store (key, value) objects for faster retrieval. 

---

### Features
* Fast storage and retrieval.
* Efficient replacement based on LRU eviction policy when trying to insert in a full cache.
* Expiry of cache entries based on TTL value specified at time of cache insert.
* Extensible to support other eviction policies.