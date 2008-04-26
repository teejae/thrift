/**
 * Thrift c
 *
 */
#ifndef THRIFT_STRUCT_H
#define THRIFT_STRUCT_H

/* base includes */
#include <glib-object.h>
#include "thrift_protocol.h"

/* struct */
typedef struct _ThriftStruct ThriftStruct;
struct _ThriftStruct
{ 
    GObject parent; 

    /* private */
    gchar * bucket;
    gchar * key;
    gchar * value;
}; 
typedef struct _ThriftStructClass ThriftStructClass;
struct _ThriftStructClass
{ 
    GObjectClass parent; 

    /* public */
    gint32 (*read) (ThriftStruct * object, ThriftProtocol * thrift_protocol);
    gint32 (*write) (ThriftStruct * object, ThriftProtocol * thrift_protocol);
}; 
GType thrift_struct_get_type (void);
void thrift_struct_instance_init (GTypeInstance *instance, gpointer g_class);
#define THRIFT_TYPE_STRUCT (thrift_struct_get_type ())
#define THRIFT_STRUCT(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), THRIFT_TYPE_STRUCT, ThriftStruct))
#define THRIFT_STRUCT_CLASS(c) (G_TYPE_CHECK_CLASS_CAST ((c), THRIFT_TYPE_STRUCT, ThriftStructClass))
#define THRIFT_IS_STRUCT(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), THRIFT_TYPE_STRUCT))
#define THRIFT_IS_STRUCT_CLASS(c) (G_TYPE_CHECK_CLASS_TYPE ((c), THRIFT_TYPE_STRUCT))
#define THRIFT_STRUCT_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), THRIFT_TYPE_STRUCT, ThriftStructClass))

gint32 thrift_struct_read (ThriftStruct * object, ThriftProtocol * thift_protocol);
gint32 thrift_struct_write (ThriftStruct * object, ThriftProtocol * thift_protocol);

#endif /* THRIFT_STRUCT_H */

