To test locally:

1. If it is needed to change amixr backend url - go to the amixr/config.go and change line 12:

   From:
   ```
   client, err := amixr.NewClient(c.Token)
   ```
   
   To:
   ```
   client, err := amixr.NewClientWithCustomUrl(c.Token, "url_you_want_to_use_locally")
   ```

2. Run (If it is needed change OS_ARCH param in GNUmakefile):
   
   ```bash
   make install
   ```

3. Import:

   ```hcl
   terraform {
      required_providers {
         amixr = {
            source = "amixr.io/alertmixer/amixr"
            version = "0.2.3"
         }
      }
   }
   ```
